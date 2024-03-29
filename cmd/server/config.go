package main

import "time"

type Config struct {
	serverAddr      string        // serverAddr store address and port to send requests to a server
	storeInterval   time.Duration // Периодичность, с которой текущие показания сервера сохраняются на диск (в секундах)
	fileStoragePath string        // Полное имя файла, куда сохраняются текущие значения
	isReqRestore    bool          // Загружать ранее сохранённые значения из файла при старте сервера
	databaseDSN     string        // Строка подключения к БД
	secretKey       string        // Ключ для подписи данных
	logLevel        string
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

func (c Config) StoreInterval() time.Duration {
	return c.storeInterval
}

func (c Config) SetStoreInterval(t time.Duration) Config {
	c.storeInterval = t
	return c
}

func (c Config) SetStoreIntervalInSeconds(s uint) Config {
	c.storeInterval = time.Duration(s) * time.Second
	return c
}

func (c Config) FileStoragePath() string {
	return c.fileStoragePath
}

func (c Config) SetFileStoragePath(path string) Config {
	c.fileStoragePath = path
	return c
}

func (c Config) IsReqRestore() bool {
	return c.isReqRestore
}

func (c Config) SetIsReqRestore(restore bool) Config {
	c.isReqRestore = restore
	return c
}

func (c Config) DatabaseDSN() string {
	return c.databaseDSN
}

func (c Config) SetDatabaseDSN(dsn string) Config {
	c.databaseDSN = dsn
	return c
}

func (c Config) SecretKey() string {
	return c.secretKey
}

func (c Config) SetSecretKey(key string) Config {
	c.secretKey = key
	return c
}
