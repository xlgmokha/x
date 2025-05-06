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
	key := x.Must(pls.GenerateRandomBytes(32))

	t.Run("Sign", func(t *testing.T) {
		data := x.Must(pls.GenerateRandomBytes(64))

		tt := []struct {
			h x.Factory[hash.Hash]
		}{
			{h: md5.New},
			{h: sha1.New},
			{h: sha256.New},
			{h: sha512.New},
		}

		for _, test := range tt {
			t.Run(fmt.Sprintf("generates an HMAC %v signature", test.h), func(t *testing.T) {
				signer := x.New[*HMACSigner](WithKey(key), WithAlgorithm(test.h))
				mac := hmac.New(test.h, key)
				mac.Write(data)
				expected := mac.Sum(nil)

				result := x.Must(signer.Sign(data))

				assert.NotEmpty(t, result)
				assert.Equal(t, expected, result)
			})
		}
	})
}
