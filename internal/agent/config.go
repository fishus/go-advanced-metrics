package agent

import "time"

type config struct {
	serverAddr     string        // serverAddr store address and port to send requests to a server
	secretKey      string        // Ключ для подписи данных
	publicKeyPath  string        // Путь до файла с публичным ключом
	logLevel       string        //
	pollInterval   time.Duration // Обновлять метрики с заданной частотой (в секундах)
	reportInterval time.Duration // Отправлять метрики на сервер с заданной частотой (в секундах)
	rateLimit      uint          // Количество одновременно исходящих запросов
}

func newConfig() config {
	return config{logLevel: "info"}
}

func (c config) ServerAddr() string {
	return c.serverAddr
}

func (c config) SetServerAddr(addr string) config {
	c.serverAddr = addr
	return c
}

func (c config) PollInterval() time.Duration {
	return c.pollInterval
}

func (c config) SetPollInterval(t time.Duration) config {
	c.pollInterval = t
	return c
}

func (c config) SetPollIntervalInSeconds(s uint) config {
	c.pollInterval = time.Duration(s) * time.Second
	return c
}

func (c config) ReportInterval() time.Duration {
	return c.reportInterval
}

func (c config) SetReportInterval(t time.Duration) config {
	c.reportInterval = t
	return c
}

func (c config) SetReportIntervalInSeconds(s uint) config {
	c.reportInterval = time.Duration(s) * time.Second
	return c
}

func (c config) SecretKey() string {
	return c.secretKey
}

func (c config) SetSecretKey(key string) config {
	c.secretKey = key
	return c
}

func (c config) PublicKeyPath() string {
	return c.publicKeyPath
}

func (c config) SetPublicKeyPath(path string) config {
	c.publicKeyPath = path
	return c
}

func (c config) RateLimit() uint {
	if c.rateLimit < 1 {
		return 1
	}
	return c.rateLimit
}

func (c config) SetRateLimit(limit uint) config {
	c.rateLimit = limit
	return c
}

func (c config) LogLevel() string {
	return c.logLevel
}
