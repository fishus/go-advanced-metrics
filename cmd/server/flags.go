package main

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v10"

	"github.com/fishus/go-advanced-metrics/internal/logger"
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

	// Флаг -i=<ЗНАЧЕНИЕ> - интервал времени в секундах, по истечении которого
	// текущие показания сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)
	storeInterval := flag.Uint("i", 300, "time interval after which the current metrics values are saved to disk (in seconds)")

	// Флаг -f=<ЗНАЧЕНИЕ> - полное имя файла, куда сохраняются текущие значения
	// (по умолчанию /tmp/metrics-db.json, пустое значение отключает функцию записи на диск).
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "full filename where the current metrics values are saved")

	// Флаг -r=true/false - загружать или нет ранее сохранённые значения из указанного файла при старте сервера (по умолчанию true).
	isReqRestore := flag.Bool("r", true, "it is required to load previously saved values from the file when the server starts")

	flag.Parse()

	return config.
		SetServerAddr(*serverAddr).
		SetStoreIntervalInSeconds(*storeInterval).
		SetFileStoragePath(*fileStoragePath).
		SetIsReqRestore(*isReqRestore)
}

func parseEnvs(config Config) Config {
	var cfg struct {
		ServerAddr      string `env:"ADDRESS"`
		StoreInterval   uint   `env:"STORE_INTERVAL"`
		FileStoragePath string `env:"FILE_STORAGE_PATH"`
		IsReqRestore    bool   `env:"RESTORE"`
	}
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Panic(err.Error(), logger.String("event", "parse env"), logger.Strings("data", os.Environ()))
	}

	if _, exists := os.LookupEnv("ADDRESS"); exists {
		config = config.SetServerAddr(cfg.ServerAddr)
	}

	if _, exists := os.LookupEnv("STORE_INTERVAL"); exists {
		config = config.SetStoreIntervalInSeconds(cfg.StoreInterval)
	}

	if _, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists {
		config = config.SetFileStoragePath(cfg.FileStoragePath)
	}

	if _, exists := os.LookupEnv("RESTORE"); exists {
		config = config.SetIsReqRestore(cfg.IsReqRestore)
	}

	return config
}
