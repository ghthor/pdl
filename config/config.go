package config

import (
	"encoding/json"
	dbconfig "github.com/ghthor/database/config"
	"io/ioutil"
)

type ServerConfig struct {
	LAddr string `json:"lAddr"`

	RedirectPort uint `json:"redirectPort"`

	SslPort uint   `json:"sslPort"`
	SslCert string `json:"sslCert"`
	SslKey  string `json:"sslKey"`

	Database dbconfig.Config `json:"database"`
}

func ReadFromFile(file string) (c ServerConfig, err error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c)
	return
}
