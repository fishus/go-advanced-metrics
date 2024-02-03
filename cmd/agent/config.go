package main

import "time"

type Config struct {
	serverAddr     string        // serverAddr store address and port to send requests to a server
	pollInterval   time.Duration // Обновлять метрики с заданной частотой (в секундах)
	reportInterval time.Duration // Отправлять метрики на сервер с заданной частотой (в секундах)
	secretKey      string        // Ключ для подписи данных
	rateLimit      uint          // Количество одновременно исходящих запросов
	logLevel       string
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

func (c Config) PollInterval() time.Duration {
	return c.pollInterval
}

func (c Config) SetPollInterval(t time.Duration) Config {
	c.pollInterval = t
	return c
}

func (c Config) SetPollIntervalInSeconds(s uint) Config {
	c.pollInterval = time.Duration(s) * time.Second
	return c
}

func (c Config) ReportInterval() time.Duration {
	return c.reportInterval
}

func (c Config) SetReportInterval(t time.Duration) Config {
	c.reportInterval = t
	return c
}

func (c Config) SetReportIntervalInSeconds(s uint) Config {
	c.reportInterval = time.Duration(s) * time.Second
	return c
}

func (c Config) SecretKey() string {
	return c.secretKey
}

func (c Config) SetSecretKey(key string) Config {
	c.secretKey = key
	return c
}

func (c Config) RateLimit() uint {
	if c.rateLimit < 1 {
		return 1
	}
	return c.rateLimit
}

func (c Config) SetRateLimit(limit uint) Config {
	c.rateLimit = limit
	return c
}
