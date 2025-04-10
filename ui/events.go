package ui

import (
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/slotopol/balance/api"
	cfg "github.com/slotopol/balance/config"
)

func (p *MainPage) OnCellSelected(id widget.TableCellID) {
	if id.Row < 0 || id.Col < 0 {
		p.OnCellUnselect()
		return
	}
	var user, ok = api.Users[cfg.UserList[id.Row]]
	if !ok {
		return
	}
	p.selIdx = id.Row
	p.selUser = user
	p.userdelBut.Enable()
	if p.admAL&api.ALbooker != 0 {
		p.walletBut.Enable()
		p.mrtpBut.Enable()
	}
	if p.admAL&api.ALadmin != 0 {
		p.accessBut.Enable()
	}
	log.Printf("selected '%s'", user.Email)
}

func (p *MainPage) OnCellUnselect() {
	p.selIdx = -1
	p.selUser = nil
	p.userdelBut.Disable()
	p.walletBut.Disable()
	p.mrtpBut.Disable()
	p.accessBut.Disable()
	log.Println("unselect user")
}

func (p *MainPage) OnUserAdd(w fyne.Window) {
	var emailEdt = widget.NewEntry()
	emailEdt.Validator = EmailValidator
	emailEdt.PlaceHolder = "test@example.com"
	var items = []*widget.FormItem{
		{Text: "Email", Widget: emailEdt, HintText: "Email of registered user"},
	}
	var dlg = dialog.NewForm("Registered email", "Add", "Cancel", items, func(b bool) {
		var err error
		if !b {
			return
		}
		var user api.User
		if user, _, err = api.ReqSignIs(emailEdt.Text); err != nil {
			log.Printf("can not detect user '%s'", emailEdt.Text)
			return
		}
		cfg.UserList = append(cfg.UserList, emailEdt.Text)
		api.Users[emailEdt.Text] = &user
		p.userTable.Refresh()
		cfg.SaveUserList()
		log.Printf("user '%s' added to list", emailEdt.Text)
	}, w)
	dlg.Resize(fyne.Size{Width: 400})
	dlg.Show()
}

func (p *MainPage) OnUserRemove(w fyne.Window) {
	var email = cfg.UserList[p.selIdx]
	var dlg = dialog.NewConfirm(
		"Confirm to remove",
		fmt.Sprintf("Confirm to remove user with email '%s' from the list. It will be removed from the list only and remains in the database.", email),
		func(confirm bool) {
			if !confirm {
				return
			}
			cfg.UserList = append(cfg.UserList[:p.selIdx], cfg.UserList[p.selIdx+1:]...)
			delete(api.Users, email)
			p.userTable.Refresh()
			cfg.SaveUserList()
			log.Printf("user '%s' removed from list", email)
		}, w)
	dlg.SetDismissText("Cancel")
	dlg.SetConfirmText("Remove")
	dlg.Show()
}

func (p *MainPage) OnUserWallet(w fyne.Window) {
	if p.selIdx < 0 || p.admAL&api.ALbooker == 0 {
		return
	}
	var userTxt = widget.NewLabel(p.selUser.Email)
	var walletEdt = widget.NewEntry()
	walletEdt.Validator = validation.NewRegexp(walletRx, "not a valid sum")
	var items = []*widget.FormItem{
		{Text: "User", Widget: userTxt, HintText: ""},
		{Text: "Sum", Widget: walletEdt, HintText: "Sum to add to user balance"},
	}
	var dlg = dialog.NewForm("Balance replenishment", "Add", "Cancel", items, func(b bool) {
		if !b {
			return
		}
		var err error
		var sum, wallet float64
		if sum, err = strconv.ParseFloat(walletEdt.Text, 64); err != nil {
			log.Printf("can not parse balance sum '%s'", walletEdt.Text)
			return
		}
		if wallet, err = api.ReqWalletAdd(p.selcid, p.selUser.UID, sum); err != nil {
			log.Printf("can not add sum '%g' to user balance", sum)
			return
		}
		var props, _ = p.selUser.GetProps(p.selcid)
		props.Wallet = wallet
		p.selUser.SetProps(p.selcid, props)
		p.userTable.Refresh()
		log.Printf("added '%g' to balance, wallet is '%g'", sum, wallet)
	}, w)
	dlg.Resize(fyne.Size{Width: 240})
	dlg.Show()
}

func (p *MainPage) OnUserMrtp(w fyne.Window) {
	if p.selIdx < 0 || p.admAL&api.ALbooker == 0 {
		return
	}
	var userTxt = widget.NewLabel(p.selUser.Email)
	var mrtpEdt = widget.NewEntry()
	mrtpEdt.Validator = MrtpValidator
	mrtpEdt.PlaceHolder = "92.5"
	var items = []*widget.FormItem{
		{Text: "User", Widget: userTxt, HintText: ""},
		{Text: "MRTP", Widget: mrtpEdt, HintText: "Master RTP percent"},
	}
	var dlg = dialog.NewForm("Master Return to Player percent", "Set", "Cancel", items, func(b bool) {
		if !b {
			return
		}
		var err error
		var mrtp float64
		if mrtp, err = strconv.ParseFloat(mrtpEdt.Text, 64); err != nil {
			log.Printf("can not parse MRTP '%s'", mrtpEdt.Text)
			return
		}
		if err = api.ReqRtpSet(p.selcid, p.selUser.UID, mrtp); err != nil {
			log.Printf("can not set MRTP '%g' to user", mrtp)
			return
		}
		var props, _ = p.selUser.GetProps(p.selcid)
		props.MRTP = mrtp
		p.selUser.SetProps(p.selcid, props)
		p.userTable.Refresh()
		log.Printf("set MRTP '%g' to user", mrtp)
	}, w)
	dlg.Resize(fyne.Size{Width: 240})
	dlg.Show()
}

func (p *MainPage) OnUserAccess(w fyne.Window) {
	if p.selIdx < 0 || p.admAL&api.ALadmin == 0 {
		return
	}
	var props, _ = p.selUser.GetProps(p.selcid)
	var access = props.Access
	var userTxt = widget.NewLabel(p.selUser.Email)
	var memberChk = widget.NewCheck("user have access to club", func(is bool) {
		if is {
			access |= api.ALmember
		} else {
			access &^= api.ALmember
		}
	})
	var dealerChk = widget.NewCheck("club game settings and users gameplay", func(is bool) {
		if is {
			access |= api.ALdealer
		} else {
			access &^= api.ALdealer
		}
	})
	var bookerChk = widget.NewCheck("user properties and manage user money", func(is bool) {
		if is {
			access |= api.ALbooker
		} else {
			access &^= api.ALbooker
		}
	})
	var masterChk = widget.NewCheck("club bank, fund, deposit", func(is bool) {
		if is {
			access |= api.ALmaster
		} else {
			access &^= api.ALmaster
		}
	})
	var adminChk = widget.NewCheck("change same access levels to others", func(is bool) {
		if is {
			access |= api.ALadmin
		} else {
			access &^= api.ALadmin
		}
	})
	memberChk.Checked = access&api.ALmember != 0
	dealerChk.Checked = access&api.ALdealer != 0
	bookerChk.Checked = access&api.ALbooker != 0
	masterChk.Checked = access&api.ALmaster != 0
	adminChk.Checked = access&api.ALadmin != 0
	if p.admAL&api.ALmember == 0 {
		memberChk.Disable()
	}
	if p.admAL&api.ALdealer == 0 {
		dealerChk.Disable()
	}
	if p.admAL&api.ALbooker == 0 {
		bookerChk.Disable()
	}
	if p.admAL&api.ALmaster == 0 {
		masterChk.Disable()
	}
	if p.admAL&api.ALadmin == 0 {
		adminChk.Disable()
	}
	var items = []*widget.FormItem{
		{Text: "User", Widget: userTxt},
		{Text: "member", Widget: memberChk},
		{Text: "dealer", Widget: dealerChk},
		{Text: "booker", Widget: bookerChk},
		{Text: "master", Widget: masterChk},
		{Text: "admin", Widget: adminChk},
	}
	var dlg = dialog.NewForm("Access rights", "Set", "Cancel", items, func(b bool) {
		if !b {
			return
		}
		var err error
		if err = api.ReqAccessSet(p.selcid, p.selUser.UID, access); err != nil {
			log.Printf("can not set access rights '%s' to user", FormatAL(access))
			return
		}
		var props, _ = p.selUser.GetProps(p.selcid)
		props.Access = access
		p.selUser.SetProps(p.selcid, props)
		p.userTable.Refresh()
		log.Printf("set access rights '%s' to user", FormatAL(access))
	}, w)
	dlg.Resize(fyne.Size{Width: 240})
	dlg.Show()
}

func (p *MainPage) OnClubBank(w fyne.Window) {
	if p.admAL&api.ALmaster == 0 {
		return
	}
	var clubname = p.clubTabs.Selected().Text
	var clubTxt = widget.NewLabel(clubname)
	var bankEdt = widget.NewEntry()
	bankEdt.Validator = validation.NewRegexp(walletRx, "not a valid sum")
	var fundEdt = widget.NewEntry()
	fundEdt.Validator = validation.NewRegexp(walletRx, "not a valid sum")
	var lockEdt = widget.NewEntry()
	lockEdt.Validator = validation.NewRegexp(walletRx, "not a valid sum")
	var items = []*widget.FormItem{
		{Text: "Club name", Widget: clubTxt, HintText: ""},
		{Text: "Bank sum", Widget: bankEdt, HintText: "Sum to add to club bank"},
		{Text: "Fund sum", Widget: fundEdt, HintText: "Sum to add to club jackpot fund"},
		{Text: "Lock sum", Widget: lockEdt, HintText: "Sum to add to club deposit"},
	}
	var dlg = dialog.NewForm("Bank transfer", "Add", "Cancel", items, func(b bool) {
		if !b {
			return
		}
		var err error
		var bsum, fsum, lsum float64
		var ret api.RetClubCashin
		if bsum, err = strconv.ParseFloat(bankEdt.Text, 64); err != nil {
			log.Printf("can not parse club bank sum '%s'", bankEdt.Text)
			return
		}
		if fsum, err = strconv.ParseFloat(fundEdt.Text, 64); err != nil {
			log.Printf("can not parse club bank sum '%s'", fundEdt.Text)
			return
		}
		if lsum, err = strconv.ParseFloat(lockEdt.Text, 64); err != nil {
			log.Printf("can not parse club bank sum '%s'", lockEdt.Text)
			return
		}
		if ret, err = api.ClubCashin(p.selcid, bsum, fsum, lsum); err != nil {
			log.Printf("can not add sum to bank '%g', jackpot fund '%g', deposit '%g'", bsum, fsum, lsum)
			return
		}
		log.Printf("add to bank '%g', to jackpot fund '%g', to deposit '%g'", bsum, fsum, lsum)
		var label = p.clubTabs.Selected().Content.(*widget.Label)
		label.SetText(fmt.Sprintf("bank: %.2f, jackpot fund: %.2f, deposit: %.2f", ret.Bank, ret.Fund, ret.Lock))
		label.Refresh()
	}, w)
	dlg.Resize(fyne.Size{Width: 240})
	dlg.Show()
}
