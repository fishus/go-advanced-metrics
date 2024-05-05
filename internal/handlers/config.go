package handlers

import (
	store "github.com/fishus/go-advanced-metrics/internal/storage"
)

var storage store.MetricsStorager

func Storage() store.MetricsStorager {
	return storage
}

func SetStorage(s store.MetricsStorager) {
	storage = s
}

var secretKey string

func SecretKey() string {
	return secretKey
}

func SetSecretKey(key string) {
	secretKey = key
}

var privateKey []byte

func PrivateKey() []byte {
	return privateKey
}

func SetPrivateKey(key []byte) {
	privateKey = key
}
