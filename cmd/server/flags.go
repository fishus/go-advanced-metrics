package main

import (
	"flag"
	"os"
)

func loadConfig() Config {
	config := NewConfig()
	config = parseFlags(config)
	config = parseEnvs(config)

	return config
}

func parseFlags(config Config) Config {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	serverAddr := flag.String("a", "localhost:8080", "address and port to run the server")
	flag.Parse()

	return config.SetServerAddr(*serverAddr)
}

func parseEnvs(config Config) Config {
	if serverAddr, exists := os.LookupEnv("ADDRESS"); exists {
		config = config.SetServerAddr(serverAddr)
	}
	return config
}
