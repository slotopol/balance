package api

import (
	"time"

	cfg "github.com/slotopol/balance/config"
)

// Access level.
type AL uint

const (
	ALmember AL = 1 << iota // user have access to club
	ALdealer                // can change club game settings
	ALbooker                // can change user balance and move user money to/from club deposit
	ALmaster                // can change club bank, fund, deposit
	ALadmin                 // can change same access levels to other users
	ALall    = ALdealer | ALbooker | ALmaster | ALadmin
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
	Users = map[string]*User{}
	Cfg   = cfg.Cfg // shortcut
)

func (p Props) Expired() bool {
	const lag = 20 * time.Millisecond // for refresh synchronization
	var d = time.Since(p.last)
	return d > Cfg.PropUpdateTick-lag
}

func (u *User) GetProps(cid uint64) (p Props, ok bool) {
	if u.props == nil {
		u.props = map[uint64]Props{} // make new empty map
		return
	}
	p, ok = u.props[cid]
	return
}

func (u *User) SetProps(cid uint64, p Props) {
	if u.props == nil {
		u.props = map[uint64]Props{} // make new empty map
	}
	u.props[cid] = p
}
