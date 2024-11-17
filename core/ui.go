package core

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	cfg "github.com/slotopol/balance/config"
)

// Foreground state of main window.
// Helps to suspend refreshes if window is on background.
var Foreground bool

var (
	curcid uint64 // current selected at tab club ID
	cural  AL     // access level of loggined account for current selected club
)

// Label compatible with ToolbarItem interface to insert into Toolbar.
type ToolbarLabel struct {
	widget.Label
}

func NewToolbarLabel(text string) *ToolbarLabel {
	var l = &ToolbarLabel{
		Label: widget.Label{
			Text:      text,
			Alignment: fyne.TextAlignLeading,
			TextStyle: fyne.TextStyle{},
		},
	}
	l.ExtendBaseWidget(l)
	return l
}

func (tl *ToolbarLabel) ToolbarObject() fyne.CanvasObject {
	tl.Label.Importance = widget.LowImportance
	return tl
}

// Layout that fits the images to whole space and cuts edges if it needs.
type FitLayout struct {
}

func (l FitLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	var ratiofit = size.Width / size.Height
	for _, child := range objects {
		var newsize = size
		var pos = fyne.NewPos(0, 0)
		if img, ok := child.(*canvas.Image); ok {
			var ratioimg = img.Aspect()
			if ratiofit > ratioimg {
				newsize.Height = size.Width / ratioimg
				pos.Y = (size.Height - newsize.Height) / 2
			} else {
				newsize.Width = size.Height * ratioimg
				pos.Y = (size.Width - newsize.Width) / 2
			}
		}
		child.Resize(newsize)
		child.Move(pos)
	}
}

func (l FitLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var minSize = fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

var underlays = []*fyne.StaticResource{
	underlay01ImgRes,
	underlay02ImgRes,
	underlay03ImgRes,
	underlay04ImgRes,
	underlay05ImgRes,
	underlay06ImgRes,
	underlay07ImgRes,
	underlay08ImgRes,
	underlay09ImgRes,
	underlay10ImgRes,
	underlay11ImgRes,
	underlay12ImgRes,
	underlay13ImgRes,
	underlay14ImgRes,
	underlay15ImgRes,
	underlay16ImgRes,
	underlay17ImgRes,
	underlay18ImgRes,
}

var (
	// Backgroud image
	underlay = &canvas.Image{
		Resource:     underlays[rand.N(len(underlays))],
		FillMode:     canvas.ImageFillContain,
		Translucency: 0.85,
	}

	// Toolbar buttons
	useraddBut = widget.NewToolbarAction(useraddIconRes, func() { fmt.Println("useradd") })
	userdelBut = widget.NewToolbarAction(userdelIconRes, func() { fmt.Println("userdel") })
	walletBut  = widget.NewToolbarAction(walletIconRes, func() { fmt.Println("wallet") })
	mrtpBut    = widget.NewToolbarAction(percentIconRes, func() { fmt.Println("mrtp") })
	accessBut  = widget.NewToolbarAction(accessIconRes, func() { fmt.Println("access") })
	bankBut    = widget.NewToolbarAction(bankIconRes, func() { fmt.Println("bank") })
	logoutBut  = widget.NewToolbarAction(logoutIconRes, func() { fmt.Println("logout") })
	loginTxt   = NewToolbarLabel("not logined yet")

	// Toolbar frame
	toolbar = widget.NewToolbar(
		useraddBut,
		userdelBut,
		widget.NewToolbarSeparator(),
		walletBut,
		mrtpBut,
		accessBut,
		widget.NewToolbarSeparator(),
		bankBut,
		widget.NewToolbarSpacer(),
		widget.NewToolbarSeparator(),
		loginTxt,
		logoutBut,
	)

	// Table with users
	clubtabs = &container.AppTabs{}
	colhdr   = []string{"email", "wallet", "MRTP", "access"}
	userlist = &widget.Table{
		Length: func() (int, int) { return len(cfg.UserList), 4 },
		CreateCell: func() fyne.CanvasObject {
			var label = widget.NewLabel("")
			label.Truncation = fyne.TextTruncateClip
			return label
		},
		UpdateCell: func(id widget.TableCellID, cell fyne.CanvasObject) {
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
		},
		UpdateHeader: func(id widget.TableCellID, cell fyne.CanvasObject) {
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
		},
		ShowHeaderRow:    true,
		ShowHeaderColumn: true,
	}

	// Main page
	mainPage = container.NewStack(
		container.New(FitLayout{}, underlay),
		container.NewBorder(
			container.NewVBox(toolbar, clubtabs),
			nil, nil, nil,
			userlist),
	)
)

// Refreshes visible content of users list. Fetches data from server
// if cached data has timeout is over.
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
