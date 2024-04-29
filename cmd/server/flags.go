package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
)

func loadConfig() Config {
	config := NewConfig()
	config = parseConfigFile(config)
	config = parseFlags(config)
	config = parseEnvs(config)

	return config
}

func parseConfigFile(config Config) Config {
	var configPath string

	// Ищем флаг -c или -config
	for i, v := range os.Args[1:] {
		switch v {
		case "-c", "-config":
			if len(os.Args) < i+3 {
				break
			}
			configPath = os.Args[i+2]
		}
	}

	// Ищем env CONFIG
	if v, exists := os.LookupEnv("CONFIG"); exists {
		configPath = v
	}

	// Загружаем переменные из конфига
	if configPath != "" {
		config = loadConfigFile(configPath, config)
	}

	return config
}

func loadConfigFile(path string, config Config) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return config
	}

	type Conf struct {
		Address       string `json:"address,omitempty"`
		ReqRestore    bool   `json:"restore,omitempty"`
		StoreInterval string `json:"store_interval,omitempty"`
		StoreFile     string `json:"store_file,omitempty"`
		DatabaseDSN   string `json:"database_dsn,omitempty"`
		CryptoKey     string `json:"crypto_key,omitempty"`
	}
	var conf Conf
	if err = json.Unmarshal(data, &conf); err != nil {
		log.Println(err)
		return config
	}

	if conf.Address != "" {
		config = config.SetServerAddr(conf.Address)
	}

	config = config.SetIsReqRestore(conf.ReqRestore)

	if conf.StoreInterval != "" {
		p, err := time.ParseDuration(conf.StoreInterval)
		if err != nil {
			log.Println(err)
		}
		config = config.SetStoreInterval(p)
	}

	if conf.StoreFile != "" {
		config = config.SetFileStoragePath(conf.StoreFile)
	}

	if conf.DatabaseDSN != "" {
		config = config.SetDatabaseDSN(conf.DatabaseDSN)
	}

	if conf.CryptoKey != "" {
		config = config.SetPrivateKeyPath(conf.CryptoKey)
	}

	return config
}

func parseFlags(config Config) Config {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	serverAddr := flag.String("a", config.serverAddr, "address and port to run the server")

	// Флаг -i=<ЗНАЧЕНИЕ> - интервал времени в секундах, по истечении которого
	// текущие показания сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)
	storeInterval := flag.Uint("i", uint(config.storeInterval.Seconds()), "time interval after which the current metrics values are saved to disk (in seconds)")

	// Флаг -f=<ЗНАЧЕНИЕ> - полное имя файла, куда сохраняются текущие значения
	// (по умолчанию /tmp/metrics-db.json, пустое значение отключает функцию записи на диск).
	fileStoragePath := flag.String("f", config.fileStoragePath, "full filename where the current metrics values are saved")

	// Флаг -r=true/false - загружать или нет ранее сохранённые значения из указанного файла при старте сервера (по умолчанию true).
	isReqRestore := flag.Bool("r", config.isReqRestore, "it is required to load previously saved values from the file when the server starts")

	// Флаг -d=<ЗНАЧЕНИЕ> - строка подключения к БД
	databaseDSN := flag.String("d", config.databaseDSN, "database URL")

	// Флаг -k=<КЛЮЧ> Ключ для подписи данных
	secretKey := flag.String("k", config.secretKey, "Secret key for signing data")

	// Флаг -crypto-key путь до файла с приватным ключом
	privateKeyPath := flag.String("crypto-key", config.privateKeyPath, "Path to the private key file")

	// Флаг -config путь к файлу конфигурации
	const configUsage = "Path to the config file"
	_ = flag.String("config", "", configUsage)
	_ = flag.String("c", "", configUsage+" (shorthand)")

	flag.Parse()

	return config.
		SetServerAddr(*serverAddr).
		SetStoreIntervalInSeconds(*storeInterval).
		SetFileStoragePath(*fileStoragePath).
		SetIsReqRestore(*isReqRestore).
		SetDatabaseDSN(*databaseDSN).
		SetSecretKey(*secretKey).
		SetPrivateKeyPath(*privateKeyPath)
}

func parseEnvs(config Config) Config {
	var cfg struct {
		ServerAddr      string `env:"ADDRESS"`
		FileStoragePath string `env:"FILE_STORAGE_PATH"`
		DatabaseDSN     string `env:"DATABASE_DSN"`
		SecretKey       string `env:"KEY"`
		PrivateKeyPath  string `env:"CRYPTO_KEY"`
		StoreInterval   uint   `env:"STORE_INTERVAL"`
		IsReqRestore    bool   `env:"RESTORE"`
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Panicln(err)
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

	if _, exists := os.LookupEnv("DATABASE_DSN"); exists {
		config = config.SetDatabaseDSN(cfg.DatabaseDSN)
	}
	if _, exists := os.LookupEnv("KEY"); exists {
		config = config.SetSecretKey(cfg.SecretKey)
	}
	if _, exists := os.LookupEnv("CRYPTO_KEY"); exists {
		config = config.SetPrivateKeyPath(cfg.PrivateKeyPath)
	}

	return config
}
