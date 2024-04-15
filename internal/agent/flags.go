package agent

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v10"

	"github.com/fishus/go-advanced-metrics/internal/logger"
)

func loadConfig() config {
	conf := newConfig()
	conf = parseFlags(conf)
	conf = parseEnvs(conf)

	return conf
}

func parseFlags(config config) config {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	serverAddr := flag.String("a", "localhost:8080", "server address and port")

	// Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	pollInterval := flag.Uint("p", 2, "frequency of polling metrics from the runtime package (in seconds)")

	// Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	reportInterval := flag.Uint("r", 10, "frequency of sending metrics to the server (in seconds)")

	// Флаг -k=<КЛЮЧ> Ключ для подписи данных
	secretKey := flag.String("k", "", "Secret key for signing data")

	// Флаг -l=<ЗНАЧЕНИЕ> Количество одновременно исходящих запросов
	rateLimit := flag.Uint("l", 3, "Количество одновременно исходящих запросов")

	flag.Parse()

	return config.
		SetServerAddr(*serverAddr).
		SetPollIntervalInSeconds(*pollInterval).
		SetReportIntervalInSeconds(*reportInterval).
		SetSecretKey(*secretKey).
		SetRateLimit(*rateLimit)
}

func parseEnvs(config config) config {
	var cfg struct {
		ServerAddr     string `env:"ADDRESS"`
		SecretKey      string `env:"KEY"`
		PollInterval   uint   `env:"POLL_INTERVAL"`
		ReportInterval uint   `env:"REPORT_INTERVAL"`
		RateLimit      uint   `env:"RATE_LIMIT"`
	}
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Panic(err.Error(), logger.String("event", "parse env"), logger.Strings("data", os.Environ()))
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

	return config
}
