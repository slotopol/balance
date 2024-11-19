package ui

import (
	"errors"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/slotopol/balance/api"
	cfg "github.com/slotopol/balance/config"
)

func (f *Frame) MakeClubList() (err error) {
	var cl api.RetClubList
	if cl, err = api.ReqClubList(); err != nil {
		return
	}

	clear(api.Clubs)
	var tabs = make([]*container.TabItem, len(cl.List))
	for i, item := range cl.List {
		api.Clubs[item.Name] = item.CID
		tabs[i] = container.NewTabItem(item.Name, widget.NewLabel(""))
	}
	f.clubTabs.SetItems(tabs)

	f.clubTabs.OnSelected = func(tab *container.TabItem) {
		var err error

		var ok bool
		if curcid, ok = api.Clubs[tab.Text]; !ok {
			return
		}
		if cural, err = api.ReqAccessGet(curcid, api.Admin.UID, true); err != nil {
			return
		}

		f.RefreshContent()
	}
	f.clubTabs.OnSelected(f.clubTabs.Selected())

	log.Printf("clubs list ready, %d clubs", len(api.Clubs))
	return
}

func (f *Frame) MakeUserList() (err error) {
	var users []api.User
	if users, err = api.ReqUserIs(cfg.UserList); err != nil {
		return
	}
	var save bool
	var i int
	for _, user := range users {
		var email = cfg.UserList[i]
		if user.UID == 0 {
			cfg.UserList = append(cfg.UserList[:i], cfg.UserList[i+1:]...)
			log.Printf("user with email '%s' presents in yaml list but absent in server database, skipped", email)
			save = true
			continue
		}
		api.Users[email] = &user
		i++
	}
	go f.RefreshLoop()
	log.Printf("users list ready, %d users", len(api.Users))
	if save {
		cfg.SaveUserList()
	}
	return
}

func (f *Frame) RefreshLoop() {
	var c = time.Tick(Cfg.PropUpdateTick)
	for range c {
		if Foreground {
			f.RefreshContent()
		}
	}
}

func WaitToken() (err error) {
	for {
		var t time.Time
		for {
			if t, err = time.Parse(api.Admin.Expire, time.RFC3339); err != nil {
				return
			}
			if !t.IsZero() {
				break
			}
			<-time.After(5 * time.Minute)
		}
		// get tokens before expire
		<-time.After(time.Until(t.Add(-15 * time.Second)))
		if api.Admin, err = api.ReqRefresh(); err != nil {
			return
		}
	}
}

func (f *Frame) StartupChain() {
	var chain = [](func() error){
		f.MakeClubList,
		f.MakeUserList,
	}
	for _, step := range chain {
		if err := step(); err != nil {
			log.Printf("startup chain does not complete: %s\n", err.Error())
			return
		}
	}
}

func (f *Frame) submitSignin() {
	if api.Admin.Email != f.email.Text {
		var err error
		if api.Admin, err = api.ReqSignIn(f.email.Text, f.secret.Text); err != nil {
			var msg string
			var aerr api.AjaxErr
			if errors.As(err, &aerr) {
				msg = aerr.What
			} else {
				msg = err.Error()
			}
			f.errmsg.SetText(fmt.Sprintf("can not sign in with given credentials, %s", msg))
			return
		}
		log.Printf("signed as '%s'", cfg.Credentials.Email)

		var save = cfg.Credentials.Addr != f.host.Text ||
			cfg.Credentials.Email != f.email.Text ||
			cfg.Credentials.Secret != f.secret.Text
		if save {
			cfg.Credentials.Addr = f.host.Text
			cfg.Credentials.Email = f.email.Text
			cfg.Credentials.Secret = f.secret.Text
			if err = cfg.SaveCredentials(); err != nil {
				log.Printf("can not save credentials: %s", err.Error())
			}
		}
	}

	f.SigninPage.form.OnCancel = func() {
		f.Window.SetContent(f.mainPage)
	}
	f.SigninPage.form.Refresh()
	f.Window.SetContent(f.mainPage)
	f.loginTxt.SetText(cfg.Credentials.Email)

	go f.StartupChain()
}

func (f *Frame) CreateWindow(a fyne.App) {
	f.MainPage.Create()
	f.SigninPage.Create()

	go WaitToken()

	var w = a.NewWindow("Balance")
	w.Resize(fyne.NewSize(540, 640))
	w.SetContent(f.signinPage)
	f.Window = w

	f.SigninPage.form.OnSubmit = f.submitSignin
	f.SigninPage.form.Refresh()
	f.userTable.SetColumnWidth(0, 180) // email
	f.userTable.SetColumnWidth(1, 100) // wallet
	f.userTable.SetColumnWidth(2, 50)  // mtrp
	f.userTable.SetColumnWidth(3, 150) // access
	f.userTable.ExtendBaseWidget(f.userTable)

	if cfg.Credentials.Addr != "" && cfg.Credentials.Email != "" {
		f.submitSignin()
	}
}
