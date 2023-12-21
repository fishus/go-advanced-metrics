package main

type Config struct {
	serverAddr string // serverAddr store address and port to send requests to a server
	logLevel   string
}

func NewConfig() Config {
	return Config{logLevel: "info"}
}

func (c Config) ServerAddr() string {
	return c.serverAddr
}

func (c Config) SetServerAddr(addr string) Config {
	c.serverAddr = addr
	return c
}
