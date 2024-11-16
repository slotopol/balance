package core

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed icon/login.svg
var loginRes []byte

//go:embed icon/logout.svg
var logoutRes []byte

//go:embed icon/person_add.svg
var useraddRes []byte

//go:embed icon/person_remove.svg
var userdelRes []byte

//go:embed icon/account_balance_wallet.svg
var walletRes []byte

//go:embed icon/percent.svg
var percentRes []byte

//go:embed icon/key.svg
var accessRes []byte

//go:embed icon/price_change.svg
var bankRes []byte

var loginIconRes = theme.NewThemedResource(&fyne.StaticResource{
	StaticName:    "login",
	StaticContent: loginRes,
})

var logoutIconRes = theme.NewThemedResource(&fyne.StaticResource{
	StaticName:    "logout",
	StaticContent: logoutRes,
})

var useraddIconRes = theme.NewThemedResource(&fyne.StaticResource{
	StaticName:    "useradd",
	StaticContent: useraddRes,
})

var userdelIconRes = theme.NewThemedResource(&fyne.StaticResource{
	StaticName:    "userdel",
	StaticContent: userdelRes,
})

var walletIconRes = theme.NewThemedResource(&fyne.StaticResource{
	StaticName:    "wallet",
	StaticContent: walletRes,
})

var percentIconRes = theme.NewThemedResource(&fyne.StaticResource{
	StaticName:    "mrtp",
	StaticContent: percentRes,
})

var accessIconRes = theme.NewThemedResource(&fyne.StaticResource{
	StaticName:    "access",
	StaticContent: accessRes,
})

var bankIconRes = theme.NewThemedResource(&fyne.StaticResource{
	StaticName:    "bank",
	StaticContent: bankRes,
})
