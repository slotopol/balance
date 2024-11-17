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

var admin AuthResp

func (f *Frame) MakeSignIn() (err error) {
	if err = cfg.ReadCredentials(); err != nil {
		log.Printf("failure on reading credentials, using default: %s\n", err.Error())
		err = nil // skip this error
		return
	}
	if admin, err = ApiSignIn(cfg.Credentials.Email, cfg.Credentials.Secret); err != nil {
		return
	}
	f.loginTxt.SetText(fmt.Sprintf(cfg.Credentials.Email))
	log.Printf("signed as '%s'", cfg.Credentials.Email)
	return
}

func (f *Frame) MakeClubList() (err error) {
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
	f.clubTabs.SetItems(tabs)

	f.clubTabs.OnSelected = func(tab *container.TabItem) {
		var err error

		var ok bool
		if curcid, ok = Clubs[tab.Text]; !ok {
			return
		}
		if cural, err = ApiAccessGet(curcid, admin.UID, true); err != nil {
			return
		}

		f.RefreshContent()
	}
	f.clubTabs.OnSelected(f.clubTabs.Selected())

	log.Printf("clubs list ready, %d clubs", len(Clubs))
	return
}

func (f *Frame) MakeUserList() (err error) {
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
	go f.RefreshLoop()
	log.Printf("users list ready, %d users", len(Users))
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
			if t, err = time.Parse(admin.Expire, time.RFC3339); err != nil {
				return
			}
			if !t.IsZero() {
				break
			}
			<-time.After(5 * time.Minute)
		}
		// get tokens before expire
		<-time.After(time.Until(t.Add(-15 * time.Second)))
		if admin, err = ApiRefresh(); err != nil {
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
}
