package api

import (
	"encoding/xml"
	"time"
)

var Admin AuthResp

func ReqPing() (err error) {
	_, _, err = HttpGet[struct{}]("/ping", "", nil)
	return
}

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

func ReqSignIs(email string) (user User, status int, err error) {
	var arg = ArgSignIs{
		Email: email,
	}
	var ret RetSignIs
	ret, status, err = HttpPost[ArgSignIs, RetSignIs]("/signis", Admin.Access, &arg)
	user.UID = ret.UID
	user.Email = ret.Email
	user.Name = ret.Name
	return
}

func ReqSignIn(email, secret string) (ret AuthResp, err error) {
	var arg ArgSignIn
	arg.Email = email
	arg.Secret = secret
	ret, _, err = HttpPost[ArgSignIn, AuthResp]("/signin", "", &arg)
	return
}

func ReqRefresh() (ret AuthResp, err error) {
	ret, _, err = HttpGet[AuthResp]("/refresh", Admin.Access, nil)
	return
}

type (
	useritem struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"user"`
		UID     uint64   `json:"uid,omitempty" yaml:"uid,omitempty" xml:"uid,attr,omitempty"`
		Email   string   `json:"email,omitempty" yaml:"email,omitempty" xml:"email,attr,omitempty"`
		Name    string   `json:"name,omitempty" yaml:"name,omitempty" xml:"name,attr,omitempty"`
	}
	ArgUserIs struct {
		XMLName xml.Name   `json:"-" yaml:"-" xml:"arg"`
		List    []useritem `json:"list" yaml:"list" xml:"list>user"`
	}
	RetUserIs struct {
		XMLName xml.Name   `json:"-" yaml:"-" xml:"ret"`
		List    []useritem `json:"list" yaml:"list" xml:"list>user"`
	}
)

func ReqUserIs(emails []string) (users []User, err error) {
	var arg ArgUserIs
	var ret RetUserIs
	arg.List = make([]useritem, len(emails))
	for i, email := range emails {
		arg.List[i].Email = email
	}
	ret, _, err = HttpPost[ArgUserIs, RetUserIs]("/user/is", Admin.Access, &arg)
	users = make([]User, len(ret.List))
	for i, item := range ret.List {
		users[i].UID = item.UID
		users[i].Email = item.Email
		users[i].Name = item.Name
	}
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
	ArgClubCashin struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		BankSum float64  `json:"banksum" yaml:"banksum" xml:"banksum"`
		FundSum float64  `json:"fundsum" yaml:"fundsum" xml:"fundsum"`
		LockSum float64  `json:"locksum" yaml:"locksum" xml:"locksum"`
	}
	RetClubCashin struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		BID     uint64   `json:"bid" yaml:"bid" xml:"bid,attr"`
		Bank    float64  `json:"bank" yaml:"bank" xml:"bank"`
		Fund    float64  `json:"fund" yaml:"fund" xml:"fund"`
		Lock    float64  `json:"lock" yaml:"lock" xml:"lock"`
	}
)

func ReqClubList() (ret RetClubList, err error) {
	ret, _, err = HttpPost[any, RetClubList]("/club/list", Admin.Access, nil)
	return
}

func ReqClubInfo(cid uint64) (ret RetClubInfo, err error) {
	var arg = ArgClubInfo{
		CID: cid,
	}
	ret, _, err = HttpPost[ArgClubInfo, RetClubInfo]("/club/info", Admin.Access, &arg)
	return
}

func ClubCashin(cid uint64, bsum, fsum, lsum float64) (ret RetClubCashin, err error) {
	var arg = ArgClubCashin{
		CID:     cid,
		BankSum: bsum,
		FundSum: fsum,
		LockSum: lsum,
	}
	ret, _, err = HttpPost[ArgClubCashin, RetClubCashin]("/club/cashin", Admin.Access, &arg)
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
	ArgWalletAdd struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr"`
		Sum     float64  `json:"sum" yaml:"sum" xml:"sum"`
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
	ArgAccessSet struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr"`
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
	ArgRtpSet struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr"`
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"`
	}
)

func ReqPropGet(cid, uid uint64) (p Props, err error) {
	var arg = ArgPropGet{
		CID: cid,
		UID: uid,
	}
	var ret RetPropGet
	ret, _, err = HttpPost[ArgPropGet, RetPropGet]("/prop/get", Admin.Access, &arg)
	p = Props{
		Wallet: ret.Wallet,
		Access: ret.Access,
		MRTP:   ret.MRTP,
		last:   time.Now(),
	}
	return
}

func ReqWalletGet(cid, uid uint64) (wallet float64, err error) {
	var arg = ArgPropGet{
		CID: cid,
		UID: uid,
	}
	var ret RetWalletGet
	ret, _, err = HttpPost[ArgPropGet, RetWalletGet]("/prop/wallet/get", Admin.Access, &arg)
	wallet = ret.Wallet
	return
}

func ReqWalletAdd(cid, uid uint64, sum float64) (wallet float64, err error) {
	var arg = ArgWalletAdd{
		CID: cid,
		UID: uid,
		Sum: sum,
	}
	var ret RetWalletGet
	ret, _, err = HttpPost[ArgWalletAdd, RetWalletGet]("/prop/wallet/add", Admin.Access, &arg)
	wallet = ret.Wallet
	return
}

func ReqAccessGet(cid, uid uint64, all bool) (al AL, err error) {
	var arg = ArgAccessGet{
		CID: cid,
		UID: uid,
		All: all,
	}
	var ret RetAccessGet
	ret, _, err = HttpPost[ArgAccessGet, RetAccessGet]("/prop/al/get", Admin.Access, &arg)
	al = ret.Access
	return
}

func ReqAccessSet(cid, uid uint64, al AL) (err error) {
	var arg = ArgAccessSet{
		CID:    cid,
		UID:    uid,
		Access: al,
	}
	_, _, err = HttpPost[ArgAccessSet, struct{}]("/prop/al/set", Admin.Access, &arg)
	return
}

func ReqRtpGet(cid, uid uint64, all bool) (mrtp float64, err error) {
	var arg = ArgRtpGet{
		CID: cid,
		UID: uid,
		All: all,
	}
	var ret RetRtpGet
	ret, _, err = HttpPost[ArgRtpGet, RetRtpGet]("/prop/rtp/get", Admin.Access, &arg)
	mrtp = ret.MRTP
	return
}

func ReqRtpSet(cid, uid uint64, mrtp float64) (err error) {
	var arg = ArgRtpSet{
		CID:  cid,
		UID:  uid,
		MRTP: mrtp,
	}
	_, _, err = HttpPost[ArgRtpSet, struct{}]("/prop/rtp/set", Admin.Access, &arg)
	return
}
