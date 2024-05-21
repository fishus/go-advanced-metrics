package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
)

func loadConfig() (conf config, err error) {
	conf = newConfig()
	conf, err = parseFlags(conf)
	if err != nil {
		return
	}

	conf, err = parseEnvs(conf)
	if err != nil {
		return
	}

	conf, err = parseConfigFile(conf)
	return
}

func parseConfigFile(config config) (config, error) {
	if config.configFile == "" {
		return config, nil
	}

	// Значения по-умолчанию
	defaults := newConfig()

	// Загружаем переменные из конфига
	cf, err := loadConfigFile(config.configFile, defaults)
	if err != nil {
		return config, err
	}

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

	if config.trustedSubnet.String() == defaults.trustedSubnet.String() && cf.trustedSubnet.String() != defaults.trustedSubnet.String() {
		config.trustedSubnet = cf.trustedSubnet
	}

	return config, nil
}

func loadConfigFile(path string, config config) (config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("can't read config file: %w", err)
	}

	type Conf struct {
		Address       string `json:"address,omitempty"`
		ReqRestore    bool   `json:"restore,omitempty"`
		StoreInterval string `json:"store_interval,omitempty"`
		StoreFile     string `json:"store_file,omitempty"`
		DatabaseDSN   string `json:"database_dsn,omitempty"`
		CryptoKey     string `json:"crypto_key,omitempty"`
		TrustedSubnet string `json:"trusted_subnet,omitempty"`
	}
	var conf Conf
	if err = json.Unmarshal(data, &conf); err != nil {
		return config, fmt.Errorf("failed to parse json data from config file: %w", err)
	}

	if conf.Address != "" {
		config = config.SetServerAddr(conf.Address)
	}

	config = config.SetIsReqRestore(conf.ReqRestore)

	if conf.StoreInterval != "" {
		p, err := time.ParseDuration(conf.StoreInterval)
		if err != nil {
			return config, fmt.Errorf("failed to parse duration in store_interval when processing config file: %w", err)
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

	if conf.TrustedSubnet != "" {
		config, err = config.SetTrustedSubnetFromString(conf.TrustedSubnet)
		if err != nil {
			return config, fmt.Errorf("failed to parse subnet in trusted_subnet when processing config file: %w", err)
		}
	}

	return config, nil
}

func parseFlags(config config) (config, error) {
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

	// Строковое представление бесклассовой адресации (CIDR).
	var t string
	if config.trustedSubnet != nil {
		t = config.trustedSubnet.String()
	}
	trustedSubnet := flag.String("t", t, "Trusted subnet (CIDR)")

	// Флаг -g запускать gRPC сервер
	useGRPC := flag.Bool("g", false, "run gRPC server instead of REST")

	// Флаг -config путь к файлу конфигурации
	const configUsage = "Path to the config file"
	flag.StringVar(&config.configFile, "config", "", configUsage)
	flag.StringVar(&config.configFile, "c", "", configUsage+" (shorthand)")

	flag.Parse()

	if *useGRPC {
		config = config.SetServerType(ServerTypeGRPC)
	}

	if *trustedSubnet != "" {
		c, err := config.SetTrustedSubnetFromString(*trustedSubnet)
		if err != nil {
			return config, fmt.Errorf("failed to parse trusted subnet: %w", err)
		}
		config = c
	}

	return config.
		SetServerAddr(*serverAddr).
		SetStoreIntervalInSeconds(*storeInterval).
		SetFileStoragePath(*fileStoragePath).
		SetIsReqRestore(*isReqRestore).
		SetDatabaseDSN(*databaseDSN).
		SetSecretKey(*secretKey).
		SetPrivateKeyPath(*privateKeyPath), nil
}

func parseEnvs(config config) (config, error) {
	var cfg struct {
		ServerAddr      string `env:"ADDRESS"`
		FileStoragePath string `env:"FILE_STORAGE_PATH"`
		DatabaseDSN     string `env:"DATABASE_DSN"`
		SecretKey       string `env:"KEY"`
		PrivateKeyPath  string `env:"CRYPTO_KEY"`
		TrustedSubnet   string `env:"TRUSTED_SUBNET"`
		ConfigFile      string `env:"CONFIG"`
		StoreInterval   uint   `env:"STORE_INTERVAL"`
		IsReqRestore    bool   `env:"RESTORE"`
	}
	err := env.Parse(&cfg)
	if err != nil {
		return config, fmt.Errorf("failed to parse environment variables: %w", err)
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

	if _, exists := os.LookupEnv("TRUSTED_SUBNET"); exists {
		c, err := config.SetTrustedSubnetFromString(cfg.TrustedSubnet)
		if err != nil {
			return config, fmt.Errorf("failed to parse trusted subnet: %w", err)
		}
		config = c
	}

	if _, exists := os.LookupEnv("CONFIG"); exists {
		config.configFile = cfg.ConfigFile
	}

	return config, nil
}
