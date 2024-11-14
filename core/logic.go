package core

import (
	"time"

	cfg "github.com/slotopol/balance/config"
)

// Access level.
type AL uint

const (
	ALmem   AL = 1 << iota // user have access to club
	ALgame                 // can change club game settings
	ALuser                 // can change user balance and move user money to/from club deposit
	ALclub                 // can change club bank, fund, deposit
	ALadmin                // can change same access levels to other users
	ALall   = ALgame | ALuser | ALclub | ALadmin
)

type (
	Props struct {
		Wallet float64 `json:"wallet" yaml:"wallet" xml:"wallet"` // in coins
		Access AL      `json:"access" yaml:"access" xml:"access"` // access level
		MRTP   float64 `json:"mrtp" yaml:"mrtp" xml:"mrtp"`       // personal master RTP
		last   time.Time
	}
	User struct {
		UID   uint64 `json:"uid" yaml:"uid" xml:"uid,attr"`                             // user ID
		Email string `json:"email" yaml:"email" xml:"email"`                            // unique user email
		Name  string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"` // user name
		props map[uint64]Props
	}
)

var (
	Clubs = map[string]uint64{}
	Users = map[string]User{}
	Cfg   = cfg.Cfg // shortcut
)

func GetProp(cid uint64, user *User) (p Props, err error) {
	const lag = 20 * time.Millisecond // for refresh synchronization
	p = user.props[curcid]
	var d = time.Since(p.last)
	if d < Cfg.PropUpdateTick-lag {
		return // return cached
	}
	if p, err = ApiPropGet(curcid, user.UID); err != nil {
		return
	}
	user.props[curcid] = p
	return
}
