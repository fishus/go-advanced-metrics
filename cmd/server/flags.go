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
	config = parseFlags(config)
	config = parseEnvs(config)
	config = parseConfigFile(config)

	return config
}

func parseConfigFile(config Config) Config {
	if config.configFile == "" {
		return config
	}

	// Значения по-умолчанию
	defaults := NewConfig()

	// Загружаем переменные из конфига
	cf := loadConfigFile(config.configFile, defaults)

	// Устанавливаем значения из файла, если они не были установлены флагом или переменной окружения

	if config.serverAddr == defaults.serverAddr && cf.serverAddr != defaults.serverAddr {
		config.serverAddr = cf.serverAddr
	}

	if config.isReqRestore == defaults.isReqRestore && cf.isReqRestore != defaults.isReqRestore {
		config.isReqRestore = cf.isReqRestore
	}

	if config.storeInterval == defaults.storeInterval && cf.storeInterval != defaults.storeInterval {
		config.storeInterval = cf.storeInterval
	}

	if config.fileStoragePath == defaults.fileStoragePath && cf.fileStoragePath != defaults.fileStoragePath {
		config.fileStoragePath = cf.fileStoragePath
	}

	if config.databaseDSN == defaults.databaseDSN && cf.databaseDSN != defaults.databaseDSN {
		config.databaseDSN = cf.databaseDSN
	}

	if config.privateKeyPath == defaults.privateKeyPath && cf.privateKeyPath != defaults.privateKeyPath {
		config.privateKeyPath = cf.privateKeyPath
	}

	return config
}

func loadConfigFile(path string, config Config) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panicln(err)
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
		log.Panicln(err)
		return config
	}

	if conf.Address != "" {
		config = config.SetServerAddr(conf.Address)
	}

	config = config.SetIsReqRestore(conf.ReqRestore)

	if conf.StoreInterval != "" {
		p, err := time.ParseDuration(conf.StoreInterval)
		if err != nil {
			log.Panicln(err)
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
	flag.StringVar(&config.configFile, "config", "", configUsage)
	flag.StringVar(&config.configFile, "c", "", configUsage+" (shorthand)")

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
		ConfigFile      string `env:"CONFIG"`
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

	if _, exists := os.LookupEnv("CONFIG"); exists {
		config.configFile = cfg.ConfigFile
	}

	return config
}
