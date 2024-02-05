package agent

var Config config

func Initialize() error {
	Config = loadConfig()
	return nil
}
