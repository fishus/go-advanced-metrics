package handlers

import (
	"net"

	"github.com/fishus/go-advanced-metrics/internal/controller"
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

var Controller controller.Controller

var config Config

type Config struct {
	ServerAddr    string
	Storage       store.MetricsStorager // TODO remove
	SecretKey     string
	PrivateKey    []byte
	TrustedSubnet *net.IPNet
}
