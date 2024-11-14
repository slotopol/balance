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

var (
	admwnd   = widget.NewLabel("not logined yet")
	clubtabs = &container.AppTabs{}
	userlist *widget.List
	curcid   uint64
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
	admwnd.SetText(fmt.Sprintf("logined as '%s'", cfg.Credentials.Email))
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
		tabs[i] = container.NewTabItem(item.Name, widget.NewLabel("Content of tab "+item.Name))
	}
	clubtabs.SetItems(tabs)
	curcid = Clubs[clubtabs.Selected().Text]

	clubtabs.OnSelected = func(tab *container.TabItem) {
		var uid, ok = Clubs[tab.Text]
		if !ok {
			return
		}
		curcid = uid
		userlist.Refresh()
	}
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
			log.Printf("user with email '%s' presents in yaml list but absent in server database, skipped")
			continue
		}
		user.props = map[uint64]Props{} // make new empty map
		Users[email] = user
	}
	go func() {
		var c = time.Tick(Cfg.PropUpdateTick)
		for range c {
			userlist.Refresh()
		}
	}()
	log.Printf("users list ready, %d users", len(Users))
	return
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

	userlist = widget.NewList(
		func() int {
			return len(cfg.UserList)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("email"), widget.NewLabelWithStyle("wallet", fyne.TextAlignLeading, fyne.TextStyle{
				Bold:      true,
				Monospace: true,
			}))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			var err error

			var mailwid = item.(*fyne.Container).Objects[0].(*widget.Label)
			mailwid.SetText(cfg.UserList[id])

			var wallwid = item.(*fyne.Container).Objects[1].(*widget.Label)
			var user, ok = Users[cfg.UserList[id]]
			if !ok {
				wallwid.SetText("error")
				return
			}
			var prop Props
			if prop, err = GetProp(curcid, &user); err != nil {
				wallwid.SetText("error")
				return
			}
			wallwid.SetText(fmt.Sprintf("%.2f", prop.Wallet))
		},
	)

	var frame = container.NewBorder(
		container.NewVBox(admwnd, clubtabs),
		nil, nil, nil,
		userlist)
	w.SetContent(frame)
	go StartupChain()
	go WaitToken()

	w.Resize(fyne.NewSize(540, 720))
	return w
}
