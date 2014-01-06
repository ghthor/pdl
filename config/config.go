package config

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {
	LAddr string `json:"lAddr"`

	RedirectPort uint `json:"redirectPort"`

	SslPort uint   `json:"sslPort"`
	SslCert string `json:"sslCert"`
	SslKey  string `json:"sslKey"`

	Database DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	DefaultDB    string `json:"defaultDB"`
	FileSystemDB string `json:"fileSystemDB"`
}

func ReadFromFile(file string) (c ServerConfig, err error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c)
	return
}
