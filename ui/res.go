package ui

import (
	_ "embed"
	"math/rand/v2"

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
	//go:embed underlay/bigcard.webp
	underlay01Res []byte

	//go:embed underlay/cards1.webp
	underlay02Res []byte

	//go:embed underlay/cards2.webp
	underlay03Res []byte

	//go:embed underlay/cards3.webp
	underlay04Res []byte

	//go:embed underlay/cards4.webp
	underlay05Res []byte

	//go:embed underlay/cards5.webp
	underlay06Res []byte

	//go:embed underlay/cards6.webp
	underlay07Res []byte

	//go:embed underlay/cards7.webp
	underlay08Res []byte

	//go:embed underlay/cards8.webp
	underlay09Res []byte

	//go:embed underlay/cards9.webp
	underlay10Res []byte

	//go:embed underlay/cards10.webp
	underlay11Res []byte

	//go:embed underlay/cards11.webp
	underlay12Res []byte

	//go:embed underlay/cards12.webp
	underlay13Res []byte

	//go:embed underlay/castle.webp
	underlay14Res []byte

	//go:embed underlay/clever.webp
	underlay15Res []byte

	//go:embed underlay/dices.webp
	underlay16Res []byte

	//go:embed underlay/dragon.webp
	underlay17Res []byte

	//go:embed underlay/scull.webp
	underlay18Res []byte
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
	underlay01ImgRes = &fyne.StaticResource{
		StaticName:    "bigcard",
		StaticContent: underlay01Res,
	}
	underlay02ImgRes = &fyne.StaticResource{
		StaticName:    "cards1",
		StaticContent: underlay02Res,
	}
	underlay03ImgRes = &fyne.StaticResource{
		StaticName:    "cards2",
		StaticContent: underlay03Res,
	}
	underlay04ImgRes = &fyne.StaticResource{
		StaticName:    "cards3",
		StaticContent: underlay04Res,
	}
	underlay05ImgRes = &fyne.StaticResource{
		StaticName:    "cards4",
		StaticContent: underlay05Res,
	}
	underlay06ImgRes = &fyne.StaticResource{
		StaticName:    "cards5",
		StaticContent: underlay06Res,
	}
	underlay07ImgRes = &fyne.StaticResource{
		StaticName:    "cards6",
		StaticContent: underlay07Res,
	}
	underlay08ImgRes = &fyne.StaticResource{
		StaticName:    "cards7",
		StaticContent: underlay08Res,
	}
	underlay09ImgRes = &fyne.StaticResource{
		StaticName:    "cards8",
		StaticContent: underlay09Res,
	}
	underlay10ImgRes = &fyne.StaticResource{
		StaticName:    "cards9",
		StaticContent: underlay10Res,
	}
	underlay11ImgRes = &fyne.StaticResource{
		StaticName:    "cards10",
		StaticContent: underlay11Res,
	}
	underlay12ImgRes = &fyne.StaticResource{
		StaticName:    "cards11",
		StaticContent: underlay12Res,
	}
	underlay13ImgRes = &fyne.StaticResource{
		StaticName:    "cards12",
		StaticContent: underlay13Res,
	}
	underlay14ImgRes = &fyne.StaticResource{
		StaticName:    "castle",
		StaticContent: underlay14Res,
	}
	underlay15ImgRes = &fyne.StaticResource{
		StaticName:    "clever",
		StaticContent: underlay15Res,
	}
	underlay16ImgRes = &fyne.StaticResource{
		StaticName:    "dices",
		StaticContent: underlay16Res,
	}
	underlay17ImgRes = &fyne.StaticResource{
		StaticName:    "dragon",
		StaticContent: underlay17Res,
	}
	underlay18ImgRes = &fyne.StaticResource{
		StaticName:    "scull",
		StaticContent: underlay18Res,
	}
)

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

func AnyUnderlay() *fyne.StaticResource {
	return underlays[rand.N(len(underlays))]
}
