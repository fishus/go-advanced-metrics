package server

import (
	"net"
	"time"
)

type config struct {
	serverAddr      string        // serverAddr store address and port to send requests to a server
	fileStoragePath string        // Полное имя файла, куда сохраняются текущие значения
	databaseDSN     string        // Строка подключения к БД
	secretKey       string        // Ключ для подписи данных
	privateKeyPath  string        // Путь до файла с приватным ключом
	logLevel        string        //
	configFile      string        // Путь к файлу конфигурации
	trustedSubnet   *net.IPNet    // Доверенная подсеть
	storeInterval   time.Duration // Периодичность, с которой текущие показания сервера сохраняются на диск (в секундах)
	isReqRestore    bool          // Загружать ранее сохранённые значения из файла при старте сервера
	serverType      ServerType
}

type ServerType string

const (
	ServerTypeREST ServerType = "rest"
	ServerTypeGRPC ServerType = "grpc"
)

func newConfig() config {
	return config{
		serverAddr:      "localhost:8080",
		fileStoragePath: "/tmp/metrics-db.json",
		logLevel:        "info",
		storeInterval:   300 * time.Second,
		isReqRestore:    true,
		serverType:      ServerTypeREST,
	}
}

func (c config) ServerAddr() string {
	return c.serverAddr
}

func (c config) SetServerAddr(addr string) config {
	c.serverAddr = addr
	return c
}

func (c config) StoreInterval() time.Duration {
	return c.storeInterval
}

func (c config) SetStoreInterval(t time.Duration) config {
	c.storeInterval = t
	return c
}

func (c config) SetStoreIntervalInSeconds(s uint) config {
	c.storeInterval = time.Duration(s) * time.Second
	return c
}

func (c config) FileStoragePath() string {
	return c.fileStoragePath
}

func (c config) SetFileStoragePath(path string) config {
	c.fileStoragePath = path
	return c
}

func (c config) IsReqRestore() bool {
	return c.isReqRestore
}

func (c config) SetIsReqRestore(restore bool) config {
	c.isReqRestore = restore
	return c
}

func (c config) DatabaseDSN() string {
	return c.databaseDSN
}

func (c config) SetDatabaseDSN(dsn string) config {
	c.databaseDSN = dsn
	return c
}

func (c config) SecretKey() string {
	return c.secretKey
}

func (c config) SetSecretKey(key string) config {
	c.secretKey = key
	return c
}

func (c config) PrivateKeyPath() string {
	return c.privateKeyPath
}

func (c config) SetPrivateKeyPath(path string) config {
	c.privateKeyPath = path
	return c
}

func (c config) TrustedSubnet() *net.IPNet {
	return c.trustedSubnet
}

func (c config) SetTrustedSubnet(subnet *net.IPNet) config {
	c.trustedSubnet = subnet
	return c
}

func (c config) SetTrustedSubnetFromString(subnet string) (config, error) {
	if subnet == "" {
		return c, nil
	}

	_, s, err := net.ParseCIDR(subnet)
	if err != nil {
		return c, err
	}

	c.trustedSubnet = s
	return c, nil
}

func (c config) ServerType() ServerType {
	return c.serverType
}

func (c config) SetServerType(t ServerType) config {
	c.serverType = t
	return c
}

func (c config) LogLevel() string {
	return c.logLevel
}
