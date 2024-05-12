package handlers

import (
	"net"

	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

var config Config

type Config struct {
	ServerAddr    string
	Storage       store.MetricsStorager
	SecretKey     string
	PrivateKey    []byte
	TrustedSubnet *net.IPNet
}
