package cfg

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	cfgcredentials = "balance-credentials.yaml"
	cfguserlist    = "balance-userlist.yaml"
)

var (
	// Executable path.
	ExePath string
	// Configuration file with path.
	CfgFile string
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

	// Config path
	const sub = "config"
	// Search config in home directory with name "balance" (without extension).
	viper.SetConfigName("balance-app")
	viper.SetConfigType("yaml")
	if env, ok := os.LookupEnv("CFGFILE"); ok {
		viper.AddConfigPath(env)
	}
	viper.AddConfigPath(filepath.Join(ExePath, sub))
	viper.AddConfigPath(ExePath)
	viper.AddConfigPath(sub)
	viper.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(home, sub))
		viper.AddConfigPath(home)
	}
	if env, ok := os.LookupEnv("GOBIN"); ok {
		viper.AddConfigPath(filepath.Join(env, sub))
		viper.AddConfigPath(env)
	} else if env, ok := os.LookupEnv("GOPATH"); ok {
		viper.AddConfigPath(filepath.Join(env, "bin", sub))
		viper.AddConfigPath(filepath.Join(env, "bin"))
	}

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		log.Println("config file not found!")
	} else {
		viper.Unmarshal(&Cfg)
		CfgFile = viper.ConfigFileUsed()
		CfgPath = filepath.Dir(CfgFile)
		log.Printf("config path: %s\n", CfgPath)
	}
}

func ReadCredentials() (err error) {
	var b []byte
	if b, err = os.ReadFile(filepath.Join(CfgPath, cfgcredentials)); err != nil {
		return
	}
	if yaml.Unmarshal(b, Credentials); err != nil {
		return
	}
	return
}

func ReadUserList() (err error) {
	var b []byte
	if b, err = os.ReadFile(filepath.Join(CfgPath, cfguserlist)); err != nil {
		return
	}
	if yaml.Unmarshal(b, &UserList); err != nil {
		return
	}
	return
}
