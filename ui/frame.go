package ui

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/slotopol/balance/api"
	cfg "github.com/slotopol/balance/config"
)

// Foreground state of main window.
// Helps to suspend refreshes if window is on background.
var Foreground bool

var (
	Cfg = cfg.Cfg // shortcut
)

const (
	hostRx   = `^((http|https|ftp):\/\/)?(\w[\w_\-]*(\.\w[\w_\-]*)*)(:\d+)?$`
	emailRx  = `^\w[\w_\-\.]*@\w+\.\w{1,4}$`
	walletRx = `^(\+|\-)?\d+(\.\d*)?$`
)

var (
	ErrBadEmail = errors.New("not a valid email")
	ErrNoUser   = errors.New("given email does not registered")
	ErrRtpMin   = errors.New("RTP is too small")
	ErrRtpMax   = errors.New("RTP is too big")
)

func GetProp(cid uint64, user *api.User) (p api.Props, err error) {
	if p, _ = user.GetProps(cid); !p.Expired() {
		return // return cached
	}
	if p, err = api.ReqPropGet(cid, user.UID); err != nil {
		return
	}
	user.SetProps(cid, p)
	return
}

func FormatAL(al api.AL) string {
	var items = make([]string, 0, 5)
	if al&api.ALmem != 0 {
		items = append(items, "member")
	}
	if al&api.ALgame != 0 {
		items = append(items, "game")
	}
	if al&api.ALuser != 0 {
		items = append(items, "user")
	}
	if al&api.ALclub != 0 {
		items = append(items, "club")
	}
	if al&api.ALadmin != 0 {
		items = append(items, "admin")
	}
	return strings.Join(items, ", ")
}

func MakeHostValidator() fyne.StringValidator {
	var rxf = validation.NewRegexp(hostRx, "not a valid host")
	return func(str string) (err error) {
		if err = rxf(str); err != nil {
			return
		}
		err = HostValidator(str)
		return
	}
}

func HostValidator(str string) error {
	var err error
	if err = api.ReqPing(); err != nil {
		return err
	}
	return nil
}

func EmailValidator(str string) error {
	var err error
	var status int
	if _, status, err = api.ReqSignIs(str); err != nil {
		switch status {
		case http.StatusBadRequest:
			return ErrBadEmail
		case http.StatusNotFound:
			return ErrNoUser
		}
		return err
	}
	return nil
}

func MrtpValidator(s string) (err error) {
	var v float64
	if v, err = strconv.ParseFloat(s, 64); err != nil {
		return
	}
	if v == 0 {
		return nil // special case
	}
	if v < 85.0 {
		return ErrRtpMin
	}
	if v > 115.0 {
		return ErrRtpMax
	}
	return
}

type Frame struct {
	fyne.Window
	SigninPage
	MainPage
}

type SigninPage struct {
	// Backgroud image
	underlay *canvas.Image

	// Form widgets
	host   *widget.Entry
	email  *widget.Entry
	secret *widget.Entry
	form   *widget.Form
	errmsg *widget.Label

	// Page frame
	signinPage *fyne.Container
}

const descrmd = `# SLOTOPOL credentials

To be able to view and change balance of users, the account must have the administrator access permission for working with *users*. To be able to view and change contents of the club bank, deposit and jackpot fund, the access permission for working with the *club* is required.
`

func (p *SigninPage) Create(w fyne.Window) {
	// Backgroud image
	p.underlay = &canvas.Image{
		Resource:     AnyUnderlay(),
		FillMode:     canvas.ImageFillContain,
		Translucency: 0.85,
	}

	// Description
	var descr = widget.NewRichTextFromMarkdown(descrmd)
	descr.Wrapping = fyne.TextWrapWord

	// Form widgets
	p.host = widget.NewEntry()
	p.host.SetPlaceHolder("http://example.com:8080")
	p.host.Validator = MakeHostValidator()
	p.host.Text = cfg.Credentials.Addr
	p.email = widget.NewEntry()
	p.email.SetPlaceHolder("test@example.com")
	p.email.Validator = EmailValidator
	p.email.Text = cfg.Credentials.Email
	p.secret = widget.NewPasswordEntry()
	p.secret.SetPlaceHolder("password")
	p.secret.Text = cfg.Credentials.Secret
	p.errmsg = widget.NewLabel("")
	p.errmsg.Wrapping = fyne.TextWrapWord
	p.form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Host", Widget: p.host, HintText: "Host address of server"},
			{Text: "Email", Widget: p.email, HintText: "A valid registered email address"},
			{Text: "Secret", Widget: p.secret, HintText: "Password for authorization"},
		},
		SubmitText: "Sign-in",
	}

	p.signinPage = container.NewStack(
		NewImageFit(p.underlay),
		container.NewVBox(
			layout.NewSpacer(),
			descr,
			p.form,
			layout.NewSpacer(),
			p.errmsg,
		),
	)
}

type MainPage struct {
	selcid  uint64    // current selected at tab club ID
	admAL   api.AL    // access level of loggined account for current selected club
	selIdx  int       // current selected row index
	selUser *api.User // user at current selected row index

	// Backgroud image
	underlay *canvas.Image

	// Toolbar buttons
	useraddBut *widget.ToolbarAction
	userdelBut *widget.ToolbarAction
	walletBut  *widget.ToolbarAction
	mrtpBut    *widget.ToolbarAction
	accessBut  *widget.ToolbarAction
	bankBut    *widget.ToolbarAction
	logoutBut  *widget.ToolbarAction
	loginTxt   *ToolbarLabel

	// Toolbar frame
	toolbar *widget.Toolbar

	// Table with users
	clubTabs  *container.AppTabs
	userTable *widget.Table
	mainPage  *fyne.Container
}

var colhdr = []string{"email", "wallet", "MRTP", "access"}

func (p *MainPage) Create(w fyne.Window) {
	// Backgroud image
	p.underlay = &canvas.Image{
		Resource:     AnyUnderlay(),
		FillMode:     canvas.ImageFillContain,
		Translucency: 0.85,
	}

	// Toolbar buttons
	p.useraddBut = widget.NewToolbarAction(useraddIconRes, func() { p.OnUserAdd(w) })
	p.userdelBut = widget.NewToolbarAction(userdelIconRes, func() { p.OnUserRemove(w) })
	p.walletBut = widget.NewToolbarAction(walletIconRes, func() { p.OnUserWallet(w) })
	p.mrtpBut = widget.NewToolbarAction(percentIconRes, func() { p.OnUserMrtp(w) })
	p.accessBut = widget.NewToolbarAction(accessIconRes, func() { p.OnUserAccess(w) })
	p.bankBut = widget.NewToolbarAction(bankIconRes, func() { fmt.Println("bank") })
	p.logoutBut = widget.NewToolbarAction(logoutIconRes, func() { fmt.Println("logout") })
	p.loginTxt = NewToolbarLabel("not logined yet")

	// Toolbar frame
	p.toolbar = widget.NewToolbar(
		p.useraddBut,
		p.userdelBut,
		widget.NewToolbarSeparator(),
		p.walletBut,
		p.mrtpBut,
		p.accessBut,
		widget.NewToolbarSeparator(),
		p.bankBut,
		widget.NewToolbarSpacer(),
		widget.NewToolbarSeparator(),
		p.loginTxt,
		p.logoutBut,
	)

	// Table with users
	p.clubTabs = &container.AppTabs{}
	p.userTable = &widget.Table{
		Length: func() (int, int) { return len(cfg.UserList), 4 },
		CreateCell: func() fyne.CanvasObject {
			var label = widget.NewLabel("")
			label.Truncation = fyne.TextTruncateClip
			return label
		},
		UpdateCell: func(id widget.TableCellID, cell fyne.CanvasObject) {
			var err error

			var label = cell.(*widget.Label)
			var user, ok = api.Users[cfg.UserList[id.Row]]
			if !ok {
				label.SetText("error")
				return
			}
			if id.Col == 0 { // email
				label.SetText(cfg.UserList[id.Row])
				return
			}
			if p.admAL&api.ALuser == 0 {
				label.SetText("N/A")
				return
			}
			var prop api.Props
			if prop, err = GetProp(p.selcid, user); err != nil {
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
		},
		UpdateHeader: func(id widget.TableCellID, cell fyne.CanvasObject) {
			var label = cell.(*widget.Label)
			if id.Row < 0 {
				label.SetText(colhdr[id.Col])
			} else if id.Col < 0 {
				var user, ok = api.Users[cfg.UserList[id.Row]]
				if !ok {
					label.SetText("error")
					return
				}
				label.SetText(strconv.Itoa(int(user.UID)))
			} else {
				label.SetText("")
			}
		},
		OnSelected:       p.OnCellSelected,
		OnUnselected:     func(id widget.TableCellID) { p.OnCellUnselect() },
		ShowHeaderRow:    true,
		ShowHeaderColumn: true,
	}

	p.OnCellUnselect()

	// Main page
	p.mainPage = container.NewStack(
		NewImageFit(p.underlay),
		container.NewBorder(
			container.NewVBox(p.toolbar, p.clubTabs),
			nil, nil, nil,
			p.userTable),
	)
}

// Refreshes visible content of users list. Fetches data from server
// if cached data has timeout is over.
func (p *MainPage) RefreshContent() {
	var err error

	p.userTable.Refresh()

	var label = p.clubTabs.Selected().Content.(*widget.Label)
	var bank, fund, deposit = "N/A", "N/A", "N/A"
	if p.admAL&api.ALclub != 0 {
		var info api.RetClubInfo
		if info, err = api.ReqClubInfo(p.selcid); err != nil {
			return
		}
		bank, fund, deposit = fmt.Sprintf("%.2f", info.Bank), fmt.Sprintf("%.2f", info.Fund), fmt.Sprintf("%.2f", info.Lock)
	}
	label.SetText(fmt.Sprintf("bank: %s, jackpot fund: %s, deposit: %s", bank, fund, deposit))
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
