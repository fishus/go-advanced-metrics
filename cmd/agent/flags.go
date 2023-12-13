package main

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
)

// serverAddr store address and port to send requests to a server
var serverAddr string

var options struct {
	pollInterval   time.Duration // Обновлять метрики с заданной частотой (в секундах)
	reportInterval time.Duration // Отправлять метрики на сервер с заданной частотой (в секундах)
}

func parseFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	flag.StringVar(&serverAddr, "a", "localhost:8080", "server address and port")

	// Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	p := flag.Uint("p", 2, "frequency of polling metrics from the runtime package (in seconds)")

	// Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	r := flag.Uint("r", 10, "frequency of sending metrics to the server (in seconds)")

	flag.Parse()

	options.pollInterval = time.Duration(*p) * time.Second
	options.reportInterval = time.Duration(*r) * time.Second

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
		serverAddr = cfg.ServerAddr
	}

	if cfg.PollInterval != 0 {
		options.pollInterval = time.Duration(cfg.PollInterval) * time.Second
	}

	if cfg.ReportInterval != 0 {
		options.reportInterval = time.Duration(cfg.ReportInterval) * time.Second
	}
}
