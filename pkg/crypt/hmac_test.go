package crypt

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xlgmokha/x/pkg/pls"
	"github.com/xlgmokha/x/pkg/x"
)

func TestHMAC(t *testing.T) {
	t.Run("Sign", func(t *testing.T) {
		data := x.Must(pls.GenerateRandomBytes(64))

		for _, hash := range []x.Factory[hash.Hash]{md5.New, sha1.New, sha256.New, sha512.New} {
			t.Run(fmt.Sprintf("generates an HMAC %v signature", hash), func(t *testing.T) {
				key := x.Must(pls.GenerateRandomBytes(32))
				signer := x.New[*HMACSigner](WithKey(key), WithAlgorithm(hash))
				mac := hmac.New(hash, key)
				mac.Write(data)
				expected := mac.Sum(nil)

				result := x.Must(signer.Sign(data))

				assert.NotEmpty(t, result)
				assert.Equal(t, expected, result)
			})
		}
	})
}
