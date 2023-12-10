package main

import (
	"flag"
	"os"
)

// serverAddr store address and port to run the server
var serverAddr string

func parseFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	flag.StringVar(&serverAddr, "a", "localhost:8080", "address and port to run the server")
	flag.Parse()
}
