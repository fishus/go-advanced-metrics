package main

import (
	"flag"
	"github.com/caarlos0/env/v10"
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
	serverAddr := flag.String("a", "localhost:8080", "server address and port")

	// Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	pollInterval := flag.Uint("p", 2, "frequency of polling metrics from the runtime package (in seconds)")

	// Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	reportInterval := flag.Uint("r", 10, "frequency of sending metrics to the server (in seconds)")

	flag.Parse()

	return config.
		SetServerAddr(*serverAddr).
		SetPollIntervalInSeconds(*pollInterval).
		SetReportIntervalInSeconds(*reportInterval)
}

func parseEnvs(config Config) Config {
	var cfg struct {
		ServerAddr     string `env:"ADDRESS"`
		PollInterval   uint   `env:"POLL_INTERVAL"`
		ReportInterval uint   `env:"REPORT_INTERVAL"`
	}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	if cfg.ServerAddr != "" {
		config = config.SetServerAddr(cfg.ServerAddr)
	}

	if cfg.PollInterval != 0 {
		config = config.SetPollIntervalInSeconds(cfg.PollInterval)
	}

	if cfg.ReportInterval != 0 {
		config = config.SetReportIntervalInSeconds(cfg.ReportInterval)
	}

	return config
}
