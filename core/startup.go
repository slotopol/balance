package core

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	cfg "github.com/slotopol/balance/config"
)

func MakeSignIn() (err error) {
	if err = cfg.ReadCredentials(); err != nil {
		log.Printf("failure on reading credentials, using default: %s\n", err.Error())
		err = nil // skip this error
		return
	}
	if admin, err = ApiSignIn(cfg.Credentials.Email, cfg.Credentials.Secret); err != nil {
		return
	}
	loginTxt.SetText(fmt.Sprintf(cfg.Credentials.Email))
	log.Printf("signed as '%s'", cfg.Credentials.Email)
	return
}

func MakeClubList() (err error) {
	var cl RetClubList
	if cl, err = ApiClubList(); err != nil {
		return
	}

	clear(Clubs)
	var tabs = make([]*container.TabItem, len(cl.List))
	for i, item := range cl.List {
		Clubs[item.Name] = item.CID
		tabs[i] = container.NewTabItem(item.Name, widget.NewLabel(""))
	}
	clubtabs.SetItems(tabs)

	clubtabs.OnSelected = func(tab *container.TabItem) {
		var err error

		var ok bool
		if curcid, ok = Clubs[tab.Text]; !ok {
			return
		}
		if cural, err = ApiAccessGet(curcid, admin.UID, true); err != nil {
			return
		}

		RefreshContent()
	}
	clubtabs.OnSelected(clubtabs.Selected())

	log.Printf("clubs list ready, %d clubs", len(Clubs))
	return
}

func MakeUserList() (err error) {
	if err = cfg.ReadUserList(); err != nil {
		log.Printf("failure on reading userlist, using default: %s\n", err.Error())
		err = nil // skip this error
	}
	for i, email := range cfg.UserList {
		var user User
		if user, err = ApiSignIs(email); err != nil {
			return
		}
		if user.UID == 0 {
			cfg.UserList = append(cfg.UserList[:i], cfg.UserList[i+1:]...)
			log.Printf("user with email '%s' presents in yaml list but absent in server database, skipped", email)
			continue
		}
		user.props = map[uint64]Props{} // make new empty map
		Users[email] = user
	}
	go RefreshLoop()
	log.Printf("users list ready, %d users", len(Users))
	return
}

func RefreshLoop() {
	var c = time.Tick(Cfg.PropUpdateTick)
	for range c {
		if Foreground {
			RefreshContent()
		}
	}
}

func WaitToken() (err error) {
	for {
		var t time.Time
		if t, err = time.Parse(admin.Expire, time.RFC3339); err != nil {
			return
		}
		// get tokens before expire
		<-time.After(time.Until(t.Add(-15 * time.Second)))
		if admin, err = ApiRefresh(); err != nil {
			return
		}
	}
}

func StartupChain() {
	var chain = [](func() error){
		MakeSignIn,
		MakeClubList,
		MakeUserList,
	}
	for _, f := range chain {
		if err := f(); err != nil {
			log.Printf("startup chain does not complete: %s\n", err.Error())
			return
		}
	}
}

func CreateMainWindow(a fyne.App) fyne.Window {
	var w = a.NewWindow("Balance")

	w.SetContent(mainPage)
	go StartupChain()
	go WaitToken()

	userlist.SetColumnWidth(0, 180) // email
	userlist.SetColumnWidth(1, 100) // wallet
	userlist.SetColumnWidth(2, 50)  // mtrp
	userlist.SetColumnWidth(3, 150) // access
	userlist.ExtendBaseWidget(userlist)
	w.Resize(fyne.NewSize(540, 640))
	return w
}
