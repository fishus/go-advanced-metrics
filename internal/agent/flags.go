package agent

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
)

func loadConfig() config {
	conf := newConfig()
	conf = parseConfigFile(conf)
	conf = parseFlags(conf)
	conf = parseEnvs(conf)

	return conf
}

func parseConfigFile(config config) config {
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

func loadConfigFile(path string, config config) config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return config
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
		log.Println(err)
		return config
	}

	if conf.Address != "" {
		config = config.SetServerAddr(conf.Address)
	}

	if conf.PollInterval != "" {
		p, err := time.ParseDuration(conf.PollInterval)
		if err != nil {
			log.Println(err)
		}
		config = config.SetPollInterval(p)
	}

	if conf.ReportInterval != "" {
		p, err := time.ParseDuration(conf.ReportInterval)
		if err != nil {
			log.Println(err)
		}
		config = config.SetReportInterval(p)
	}

	if conf.RateLimit != 0 {
		config = config.SetRateLimit(conf.RateLimit)
	}

	if conf.CryptoKey != "" {
		config = config.SetPublicKeyPath(conf.CryptoKey)
	}

	return config
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

	// Флаг -config путь к файлу конфигурации
	const configUsage = "Path to the config file"
	_ = flag.String("config", "", configUsage)
	_ = flag.String("c", "", configUsage+" (shorthand)")

	flag.Parse()

	return config.
		SetServerAddr(*serverAddr).
		SetPollIntervalInSeconds(*pollInterval).
		SetReportIntervalInSeconds(*reportInterval).
		SetSecretKey(*secretKey).
		SetPublicKeyPath(*publicKeyPath).
		SetRateLimit(*rateLimit)
}

func parseEnvs(config config) config {
	var cfg struct {
		ServerAddr     string `env:"ADDRESS"`
		SecretKey      string `env:"KEY"`
		PublicKeyPath  string `env:"CRYPTO_KEY"`
		PollInterval   uint   `env:"POLL_INTERVAL"`
		ReportInterval uint   `env:"REPORT_INTERVAL"`
		RateLimit      uint   `env:"RATE_LIMIT"`
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Panicln(err)
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

	return config
}
