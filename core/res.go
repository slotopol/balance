package core

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	_ "golang.org/x/image/webp"
)

// icons
var (
	//go:embed icon/login.svg
	loginRes []byte

	//go:embed icon/logout.svg
	logoutRes []byte

	//go:embed icon/person_add.svg
	useraddRes []byte

	//go:embed icon/person_remove.svg
	userdelRes []byte

	//go:embed icon/account_balance_wallet.svg
	walletRes []byte

	//go:embed icon/percent.svg
	percentRes []byte

	//go:embed icon/key.svg
	accessRes []byte

	//go:embed icon/price_change.svg
	bankRes []byte
)

// underlays
var (
	//go:embed underlay/cards1.webp
	underlay1Res []byte

	//go:embed underlay/cards2.webp
	underlay2Res []byte

	//go:embed underlay/cards3.webp
	underlay3Res []byte

	//go:embed underlay/cards4.webp
	underlay4Res []byte

	//go:embed underlay/cards5.webp
	underlay5Res []byte

	//go:embed underlay/cards6.webp
	underlay6Res []byte

	//go:embed underlay/castle.webp
	underlay7Res []byte

	//go:embed underlay/clever.webp
	underlay8Res []byte

	//go:embed underlay/dices.webp
	underlay9Res []byte

	//go:embed underlay/dragon.webp
	underlay10Res []byte
)

// icon resources
var (
	loginIconRes = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "login",
		StaticContent: loginRes,
	})
	logoutIconRes = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "logout",
		StaticContent: logoutRes,
	})
	useraddIconRes = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "useradd",
		StaticContent: useraddRes,
	})
	userdelIconRes = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "userdel",
		StaticContent: userdelRes,
	})
	walletIconRes = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "wallet",
		StaticContent: walletRes,
	})
	percentIconRes = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "mrtp",
		StaticContent: percentRes,
	})
	accessIconRes = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "access",
		StaticContent: accessRes,
	})
	bankIconRes = theme.NewThemedResource(&fyne.StaticResource{
		StaticName:    "bank",
		StaticContent: bankRes,
	})
)

// underlay resources
var (
	underlay1ImgRes = &fyne.StaticResource{
		StaticName:    "cards1",
		StaticContent: underlay1Res,
	}
	underlay2ImgRes = &fyne.StaticResource{
		StaticName:    "cards2",
		StaticContent: underlay2Res,
	}
	underlay3ImgRes = &fyne.StaticResource{
		StaticName:    "cards3",
		StaticContent: underlay3Res,
	}
	underlay4ImgRes = &fyne.StaticResource{
		StaticName:    "cards4",
		StaticContent: underlay4Res,
	}
	underlay5ImgRes = &fyne.StaticResource{
		StaticName:    "cards5",
		StaticContent: underlay5Res,
	}
	underlay6ImgRes = &fyne.StaticResource{
		StaticName:    "cards6",
		StaticContent: underlay6Res,
	}
	underlay7ImgRes = &fyne.StaticResource{
		StaticName:    "castle",
		StaticContent: underlay7Res,
	}
	underlay8ImgRes = &fyne.StaticResource{
		StaticName:    "clever",
		StaticContent: underlay8Res,
	}
	underlay9ImgRes = &fyne.StaticResource{
		StaticName:    "dices",
		StaticContent: underlay9Res,
	}
	underlay10ImgRes = &fyne.StaticResource{
		StaticName:    "dragon",
		StaticContent: underlay10Res,
	}
)
