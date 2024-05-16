package agent

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
	conf = parseFlags(conf)
	conf, err = parseEnvs(conf)
	if err != nil {
		return
	}

	conf, err = parseConfigFile(conf)
	if err != nil {
		return
	}

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

	if config.publicKeyPath == defaults.publicKeyPath && cf.publicKeyPath != defaults.publicKeyPath {
		config.publicKeyPath = cf.publicKeyPath
	}

	if config.pollInterval == defaults.pollInterval && cf.pollInterval != defaults.pollInterval {
		config.pollInterval = cf.pollInterval
	}

	if config.reportInterval == defaults.reportInterval && cf.reportInterval != defaults.reportInterval {
		config.reportInterval = cf.reportInterval
	}

	if config.rateLimit == defaults.rateLimit && cf.rateLimit != defaults.rateLimit {
		config.rateLimit = cf.rateLimit
	}

	return config, nil
}

func loadConfigFile(path string, config config) (config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("can't read config file: %w", err)
	}

	type Conf struct {
		Address        string `json:"address,omitempty"`
		PollInterval   string `json:"poll_interval,omitempty"`
		ReportInterval string `json:"report_interval,omitempty"`
		RateLimit      uint   `json:"rate_limit,omitempty"`
		CryptoKey      string `json:"crypto_key,omitempty"`
	}
	var conf Conf
	if err = json.Unmarshal(data, &conf); err != nil {
		return config, fmt.Errorf("failed to parse json data from config file: %w", err)
	}

	if conf.Address != "" {
		config = config.SetServerAddr(conf.Address)
	}

	if conf.PollInterval != "" {
		p, err := time.ParseDuration(conf.PollInterval)
		if err != nil {
			return config, fmt.Errorf("failed to parse duration in poll_interval when processing config file: %w", err)
		}
		config = config.SetPollInterval(p)
	}

	if conf.ReportInterval != "" {
		p, err := time.ParseDuration(conf.ReportInterval)
		if err != nil {
			return config, fmt.Errorf("failed to parse duration in report_interval when processing config file: %w", err)
		}
		config = config.SetReportInterval(p)
	}

	if conf.RateLimit != 0 {
		config = config.SetRateLimit(conf.RateLimit)
	}

	if conf.CryptoKey != "" {
		config = config.SetPublicKeyPath(conf.CryptoKey)
	}

	return config, nil
}

func parseFlags(config config) config {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	serverAddr := flag.String("a", config.serverAddr, "server address and port")

	// Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	pollInterval := flag.Uint("p", uint(config.pollInterval.Seconds()), "frequency of polling metrics from the runtime package (in seconds)")

	// Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	reportInterval := flag.Uint("r", uint(config.reportInterval.Seconds()), "frequency of sending metrics to the server (in seconds)")

	// Флаг -k=<КЛЮЧ> Ключ для подписи данных
	secretKey := flag.String("k", config.secretKey, "Secret key for signing data")

	// Флаг -l=<ЗНАЧЕНИЕ> Количество одновременно исходящих запросов
	rateLimit := flag.Uint("l", config.rateLimit, "Количество одновременно исходящих запросов")

	// Флаг -crypto-key путь до файла с публичным ключом
	publicKeyPath := flag.String("crypto-key", config.publicKeyPath, "Path to the public key file")

	// Флаг -g запускать gRPC сервер
	useGRPC := flag.Bool("g", false, "run gRPC server instead of REST")

	// Флаг -config путь к файлу конфигурации
	const configUsage = "Path to the config file"
	flag.StringVar(&config.configFile, "config", "", configUsage)
	flag.StringVar(&config.configFile, "c", "", configUsage+" (shorthand)")

	flag.Parse()

	if *useGRPC {
		config = config.SetClientType(ClientTypeGRPC)
	}

	return config.
		SetServerAddr(*serverAddr).
		SetPollIntervalInSeconds(*pollInterval).
		SetReportIntervalInSeconds(*reportInterval).
		SetSecretKey(*secretKey).
		SetPublicKeyPath(*publicKeyPath).
		SetRateLimit(*rateLimit)
}

func parseEnvs(config config) (config, error) {
	var cfg struct {
		ServerAddr     string `env:"ADDRESS"`
		SecretKey      string `env:"KEY"`
		PublicKeyPath  string `env:"CRYPTO_KEY"`
		ConfigFile     string `env:"CONFIG"`
		PollInterval   uint   `env:"POLL_INTERVAL"`
		ReportInterval uint   `env:"REPORT_INTERVAL"`
		RateLimit      uint   `env:"RATE_LIMIT"`
	}
	err := env.Parse(&cfg)
	if err != nil {
		return config, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	if _, exists := os.LookupEnv("ADDRESS"); exists {
		config = config.SetServerAddr(cfg.ServerAddr)
	}

	if _, exists := os.LookupEnv("POLL_INTERVAL"); exists {
		config = config.SetPollIntervalInSeconds(cfg.PollInterval)
	}

	if _, exists := os.LookupEnv("REPORT_INTERVAL"); exists {
		config = config.SetReportIntervalInSeconds(cfg.ReportInterval)
	}

	if _, exists := os.LookupEnv("KEY"); exists {
		config = config.SetSecretKey(cfg.SecretKey)
	}

	if _, exists := os.LookupEnv("RATE_LIMIT"); exists {
		config = config.SetRateLimit(cfg.RateLimit)
	}

	if _, exists := os.LookupEnv("CRYPTO_KEY"); exists {
		config = config.SetPublicKeyPath(cfg.PublicKeyPath)
	}

	if _, exists := os.LookupEnv("CONFIG"); exists {
		config.configFile = cfg.ConfigFile
	}

	return config, nil
}
