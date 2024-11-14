package cfg

import "time"

var (
	// Admin credentials.
	Credentials = &struct {
		Addr   string `json:"hostaddr" yaml:"hostaddr" mapstructure:"hostaddr"`
		Email  string `json:"email" yaml:"email" mapstructure:"email"`
		Secret string `json:"secret" yaml:"secret" mapstructure:"secret"`
	}{
		Addr:   "http://localhost:8080",
		Email:  "admin@example.org",
		Secret: "0YBoaT",
	}

	// List of registered users emails.
	UserList = []string{
		"admin@example.org",
		"dealer@example.org",
		"player@example.org",
	}

	// Config is common application settings.
	Cfg = &struct {
		PropUpdateTick time.Duration `json:"prop-update-tick" yaml:"prop-update-tick" mapstructure:"prop-update-tick"`
	}{
		PropUpdateTick: time.Millisecond * 4000,
	}
)
