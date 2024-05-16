package agent

import (
	"github.com/fishus/go-advanced-metrics/internal/cryptokey"
)

var (
	Config    config
	PublicKey []byte
)

func Initialize() error {
	c, err := loadConfig()
	if err != nil {
		return err
	}

	Config = c

	if Config.publicKeyPath != "" {
		pubKey, err := cryptokey.ReadKeyFile(Config.publicKeyPath)
		if err != nil {
			return err
		}
		PublicKey = pubKey
	}

	return nil
}
