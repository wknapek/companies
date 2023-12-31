package config

import (
	"encoding/json"
	"os"
)

type AppConfiguration struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	DBURI       string `json:"dburi"`
	SessionTime int    `json:"sessionTime"`
	Port        string `json:"port"`
	DBUser      string `json:"db_user"`
	DBPasswd    string `json:"db_passwd"`
	JWT_Key     string `json:"JWT_Key"`
}

// ConfFile - default conf file name
const ConfFile = "resource/config.json"

// ReadConfiguration -
func ReadConfiguration(path string) (AppConfiguration, error) {
	var cfg AppConfiguration
	//nolint:gosec
	data, err := os.ReadFile(path)
	if err != nil {
		return AppConfiguration{}, err
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return AppConfiguration{}, err
	}

	return cfg, nil
}
