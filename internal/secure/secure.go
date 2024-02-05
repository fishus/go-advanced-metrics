package secure

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
)

func Hash(data, key []byte) []byte {
	sign := NewSign(key)
	sign.Write(data)
	return sign.Sum()
}

type Sign struct {
	key  []byte
	hash hash.Hash
}

func (s *Sign) SetKey(key []byte) {
	s.key = key
}

func (s *Sign) Write(b []byte) (n int, err error) {
	return s.hash.Write(b)
}

func (s *Sign) Sum() []byte {
	return s.hash.Sum(nil)
}

func NewSign(key []byte) *Sign {
	return &Sign{
		key:  key,
		hash: hmac.New(sha256.New, key),
	}
}
