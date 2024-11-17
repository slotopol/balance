package core

import (
	"encoding/xml"
	"time"
)

type (
	ArgSignIs struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid"`
		Email   string   `json:"email" yaml:"email" xml:"email"`
	}
	RetSignIs struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid"`
		Email   string   `json:"email" yaml:"email" xml:"email"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	}
	AuthResp struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid"`
		Email   string   `json:"email" yaml:"email" xml:"email"`
		Access  string   `json:"access" yaml:"access" xml:"access"`
		Refrsh  string   `json:"refrsh" yaml:"refrsh" xml:"refrsh"`
		Expire  string   `json:"expire" yaml:"expire" xml:"expire"`
		Living  string   `json:"living" yaml:"living" xml:"living"`
	}
	ArgSignIn struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid"`
		Email   string   `json:"email" yaml:"email" xml:"email"`
		Secret  string   `json:"secret" yaml:"secret,omitempty" xml:"secret,omitempty"`
		HS256   string   `json:"hs256,omitempty" yaml:"hs256,omitempty" xml:"hs256,omitempty"`
		SigTime string   `json:"sigtime,omitempty" yaml:"sigtime,omitempty" xml:"sigtime,omitempty"`
		Code    uint32   `json:"code,omitempty" yaml:"code,omitempty" xml:"code,omitempty"`
	}
)

func ApiSignIs(email string) (user User, err error) {
	var arg = ArgSignIs{
		Email: email,
	}
	var ret RetSignIs
	ret, _, err = HttpPost[ArgSignIs, RetSignIs]("/user/is", admin.Access, &arg)
	user.UID = ret.UID
	user.Email = ret.Email
	user.Name = ret.Name
	return
}

func ApiSignIn(email, secret string) (ret AuthResp, err error) {
	var arg ArgSignIn
	arg.Email = email
	arg.Secret = secret
	ret, _, err = HttpPost[ArgSignIn, AuthResp]("/signin", "", &arg)
	return
}

func ApiRefresh() (ret AuthResp, err error) {
	ret, _, err = HttpPost[any, AuthResp]("/refresh", admin.Access, nil)
	return
}

type (
	clubitem struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"club"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	}
	RetClubList struct {
		XMLName xml.Name   `json:"-" yaml:"-" xml:"ret"`
		List    []clubitem `json:"list" yaml:"list" xml:"list>club" form:"list"`
	}
	ArgClubInfo struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
	}
	RetClubInfo struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
		Bank    float64  `json:"bank" yaml:"bank" xml:"bank"` // users win/lost balance, in coins
		Fund    float64  `json:"fund" yaml:"fund" xml:"fund"` // jackpot fund, in coins
		Lock    float64  `json:"lock" yaml:"lock" xml:"lock"` // not changed deposit within games
		Rate    float64  `json:"rate" yaml:"rate" xml:"rate"` // jackpot rate for games with progressive jackpot
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"` // master RTP
	}
)

func ApiClubList() (ret RetClubList, err error) {
	ret, _, err = HttpPost[any, RetClubList]("/club/list", admin.Access, nil)
	return
}

func ApiClubInfo(cid uint64) (ret RetClubInfo, err error) {
	var arg = ArgClubInfo{
		CID: cid,
	}
	ret, _, err = HttpPost[ArgClubInfo, RetClubInfo]("/club/info", admin.Access, &arg)
	return
}

type (
	ArgPropGet struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr"`
	}
	RetPropGet struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
		Access  AL       `json:"access" yaml:"access" xml:"access"`
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"`
	}
	RetWalletGet struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}
	ArgAccessGet struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr"`
		All     bool     `json:"all" yaml:"all" xml:"all,attr"`
	}
	RetAccessGet struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Access  AL       `json:"access" yaml:"access" xml:"access"`
	}
	ArgRtpGet struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr"`
		All     bool     `json:"all" yaml:"all" xml:"all,attr"`
	}
	RetRtpGet struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"`
	}
)

func ApiPropGet(cid, uid uint64) (p Props, err error) {
	var arg = ArgPropGet{
		CID: cid,
		UID: uid,
	}
	var ret RetPropGet
	ret, _, err = HttpPost[ArgPropGet, RetPropGet]("/prop/get", admin.Access, &arg)
	p = Props{
		Wallet: ret.Wallet,
		Access: ret.Access,
		MRTP:   ret.MRTP,
		last:   time.Now(),
	}
	return
}

func ApiWalletGet(cid, uid uint64) (sum float64, err error) {
	var arg = ArgPropGet{
		CID: cid,
		UID: uid,
	}
	var ret RetWalletGet
	ret, _, err = HttpPost[ArgPropGet, RetWalletGet]("/prop/wallet/get", admin.Access, &arg)
	sum = ret.Wallet
	return
}

func ApiAccessGet(cid, uid uint64, all bool) (al AL, err error) {
	var arg = ArgAccessGet{
		CID: cid,
		UID: uid,
		All: all,
	}
	var ret RetAccessGet
	ret, _, err = HttpPost[ArgAccessGet, RetAccessGet]("/prop/al/get", admin.Access, &arg)
	al = ret.Access
	return
}

func ApiRtpGet(cid, uid uint64, all bool) (mrtp float64, err error) {
	var arg = ArgRtpGet{
		CID: cid,
		UID: uid,
		All: all,
	}
	var ret RetRtpGet
	ret, _, err = HttpPost[ArgRtpGet, RetRtpGet]("/prop/rtp/get", admin.Access, &arg)
	mrtp = ret.MRTP
	return
}
