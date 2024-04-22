package agent

import "github.com/fishus/go-advanced-metrics/internal/cryptokey"

var Config config

var publicKey []byte

func Initialize() error {
	Config = loadConfig()
	if Config.publicKeyPath != "" {
		pubKey, err := cryptokey.ReadKeyFile(Config.publicKeyPath)
		if err != nil {
			return err
		}
		publicKey = pubKey
	}
	return nil
}
