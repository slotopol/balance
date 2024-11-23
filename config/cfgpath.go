package cfg

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const AppID = "slotopol.balance"

const (
	cfgcredentials = "credentials.yaml"
	cfguserlist    = "userlist.yaml"
)

var (
	// Executable path.
	ExePath string
	// Configuration path.
	CfgPath string
)

func init() {
	var err error

	// Executable path
	ExePath = func() string {
		if str, err := os.Executable(); err == nil {
			return filepath.Dir(str)
		} else {
			return filepath.Dir(os.Args[0])
		}
	}()

	// Configuration path
	var oscfgpath string
	if oscfgpath, err = os.UserConfigDir(); err != nil {
		log.Printf("can not obtain user config directory, any settings can not be saved: %s", err.Error())
		return
	}
	var appdata = filepath.Join(oscfgpath, "fyne", AppID)
	if err = os.MkdirAll(appdata, 0750); err != nil {
		log.Printf("cannot create application config directory: %s\n", err.Error())
		return
	}
	CfgPath = appdata
	log.Printf("config path: %s\n", CfgPath)

	if err = ReadCredentials(); err != nil {
		log.Printf("failure on reading credentials, using default: %s\n", err.Error())
	}
	if err = ReadUserList(); err != nil {
		log.Printf("failure on reading userlist, using default: %s\n", err.Error())
	}
}

var (
	ErrNoCfgPath = errors.New("configuration path does not obtained")
)

func ReadCredentials() (err error) {
	if CfgPath == "" {
		return ErrNoCfgPath
	}
	var b []byte
	if b, err = os.ReadFile(filepath.Join(CfgPath, cfgcredentials)); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, Credentials); err != nil {
		return
	}
	return
}

func SaveCredentials() (err error) {
	defer func() {
		if err != nil {
			log.Printf("can not save credentials: %s", err.Error())
		}
	}()
	if CfgPath == "" {
		return ErrNoCfgPath
	}
	var b []byte
	if b, err = yaml.Marshal(Credentials); err != nil {
		return
	}
	if err = os.WriteFile(filepath.Join(CfgPath, cfgcredentials), b, 0666); err != nil {
		return
	}
	return
}

func ReadUserList() (err error) {
	if CfgPath == "" {
		return ErrNoCfgPath
	}
	var b []byte
	if b, err = os.ReadFile(filepath.Join(CfgPath, cfguserlist)); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &UserList); err != nil {
		return
	}
	return
}

func SaveUserList() (err error) {
	defer func() {
		if err != nil {
			log.Printf("can not save user list: %s", err.Error())
		}
	}()
	if CfgPath == "" {
		return ErrNoCfgPath
	}
	var b []byte
	if b, err = yaml.Marshal(&UserList); err != nil {
		return
	}
	if err = os.WriteFile(filepath.Join(CfgPath, cfguserlist), b, 0666); err != nil {
		return
	}
	return
}
