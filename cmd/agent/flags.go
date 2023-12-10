package main

import (
	"flag"
	"os"
	"time"
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
	p := flag.Int64("p", 2, "frequency of polling metrics from the runtime package (in seconds)")

	// Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	r := flag.Int64("r", 10, "frequency of sending metrics to the server (in seconds)")

	flag.Parse()

	options.pollInterval = time.Duration(*p) * time.Second
	options.reportInterval = time.Duration(*r) * time.Second
}
