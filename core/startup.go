package core

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	cfg "github.com/slotopol/balance/config"
)

var Foreground bool

var (
	admwnd   = widget.NewLabel("not logined yet")
	clubtabs = &container.AppTabs{}
	userlist *widget.Table
	curcid   uint64
	cural    AL
)

var colhdr = []string{
	"email", "wallet", "MRTP", "access",
}

func RefreshContent() {
	var err error

	userlist.Refresh()

	var label = clubtabs.Selected().Content.(*widget.Label)
	var bank, fund, deposit = "N/A", "N/A", "N/A"
	if cural&ALclub != 0 {
		var info RetClubInfo
		if info, err = ApiClubInfo(curcid); err != nil {
			return
		}
		bank, fund, deposit = fmt.Sprintf("%.2f", info.Bank), fmt.Sprintf("%.2f", info.Fund), fmt.Sprintf("%.2f", info.Lock)
	}
	label.SetText(fmt.Sprintf("bank: %s, jackpot fund: %s, deposit: %s", bank, fund, deposit))
}

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
	go func() {
		var c = time.Tick(Cfg.PropUpdateTick)
		for range c {
			if Foreground {
				RefreshContent()
			}
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

func Lifecycle(a fyne.App) {
	var l = a.Lifecycle()
	l.SetOnStarted(func() {
		log.Println("lifecycle: started")
	})
	l.SetOnStopped(func() {
		log.Println("lifecycle: stopped")
	})
	l.SetOnEnteredForeground(func() {
		Foreground = true
		log.Println("lifecycle: entered foreground")
	})
	l.SetOnExitedForeground(func() {
		Foreground = false
		log.Println("lifecycle: exited foreground")
	})
}

func CreateMainWindow(a fyne.App) fyne.Window {
	var w = a.NewWindow("Balance")

	userlist = widget.NewTableWithHeaders(
		func() (int, int) { return len(cfg.UserList), 4 },
		func() fyne.CanvasObject {
			var label = widget.NewLabel("")
			label.Truncation = fyne.TextTruncateClip
			return label
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			var err error

			var label = cell.(*widget.Label)
			var user, ok = Users[cfg.UserList[id.Row]]
			if !ok {
				label.SetText("error")
				return
			}
			if id.Col == 0 { // email
				label.SetText(cfg.UserList[id.Row])
				return
			}
			if cural&ALuser == 0 {
				label.SetText("N/A")
				return
			}
			var prop Props
			if prop, err = GetProp(curcid, &user); err != nil {
				label.SetText("error")
				return
			}
			switch id.Col {
			case 1: // wallet
				label.SetText(fmt.Sprintf("%.2f", prop.Wallet))
			case 2: // mtrp
				if prop.MRTP > 0 {
					label.SetText(fmt.Sprintf("%g%%", prop.MRTP))
				} else {
					label.SetText("-")
				}
			case 3: // access
				label.SetText(FormatAL(prop.Access))
			}
		})
	userlist.SetColumnWidth(0, 180) // email
	userlist.SetColumnWidth(1, 100) // wallet
	userlist.SetColumnWidth(2, 50)  // mtrp
	userlist.SetColumnWidth(3, 150) // access
	userlist.UpdateHeader = func(id widget.TableCellID, cell fyne.CanvasObject) {
		var label = cell.(*widget.Label)
		if id.Row < 0 {
			label.SetText(colhdr[id.Col])
		} else if id.Col < 0 {
			var user, ok = Users[cfg.UserList[id.Row]]
			if !ok {
				label.SetText("error")
				return
			}
			label.SetText(strconv.Itoa(int(user.UID)))
		} else {
			label.SetText("")
		}
	}

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
