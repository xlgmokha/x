package pls

import (
	"crypto/hmac"
	"crypto/sha256"
)

func GenerateHMAC256(key, value []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(value)
	return mac.Sum(nil)
}
