package server

import (
	"github.com/fishus/go-advanced-metrics/internal/cryptokey"
)

var Config config

var PrivateKey []byte

func Initialize() error {
	c, err := loadConfig()
	if err != nil {
		return err
	}

	Config = c

	if Config.privateKeyPath != "" {
		privKey, err := cryptokey.ReadKeyFile(Config.privateKeyPath)
		if err != nil {
			return err
		}
		PrivateKey = privKey
	}

	return nil
}
