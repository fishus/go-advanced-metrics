package secure

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Hash(src, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(src)
	return h.Sum(nil)
}
