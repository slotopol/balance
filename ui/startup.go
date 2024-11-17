package ui

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/slotopol/balance/api"
	cfg "github.com/slotopol/balance/config"
)

func (f *Frame) MakeSignIn() (err error) {
	if err = cfg.ReadCredentials(); err != nil {
		log.Printf("failure on reading credentials, using default: %s\n", err.Error())
		err = nil // skip this error
		return
	}
	if api.Admin, err = api.ReqSignIn(cfg.Credentials.Email, cfg.Credentials.Secret); err != nil {
		return
	}
	f.loginTxt.SetText(fmt.Sprintf(cfg.Credentials.Email))
	log.Printf("signed as '%s'", cfg.Credentials.Email)
	return
}

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
	if err = cfg.ReadUserList(); err != nil {
		log.Printf("failure on reading userlist, using default: %s\n", err.Error())
		err = nil // skip this error
	}
	for i, email := range cfg.UserList {
		var user api.User
		if user, err = api.ReqSignIs(email); err != nil {
			return
		}
		if user.UID == 0 {
			cfg.UserList = append(cfg.UserList[:i], cfg.UserList[i+1:]...)
			log.Printf("user with email '%s' presents in yaml list but absent in server database, skipped", email)
			continue
		}
		api.Users[email] = &user
	}
	go f.RefreshLoop()
	log.Printf("users list ready, %d users", len(api.Users))
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
		f.MakeSignIn,
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

func (f *Frame) CreateWindow(a fyne.App) {
	f.MainPage.Create()

	go f.StartupChain()
	go WaitToken()

	var w = a.NewWindow("Balance")
	w.Resize(fyne.NewSize(540, 640))
	w.SetContent(f.mainPage)
	f.Window = w

	f.userTable.SetColumnWidth(0, 180) // email
	f.userTable.SetColumnWidth(1, 100) // wallet
	f.userTable.SetColumnWidth(2, 50)  // mtrp
	f.userTable.SetColumnWidth(3, 150) // access
	f.userTable.ExtendBaseWidget(f.userTable)
}
