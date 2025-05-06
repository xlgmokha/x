package crypt

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xlgmokha/x/pkg/pls"
	"github.com/xlgmokha/x/pkg/x"
)

func TestHMAC(t *testing.T) {
	data := x.Must(pls.GenerateRandomBytes(64))

	for _, hash := range []x.Factory[hash.Hash]{md5.New, sha1.New, sha256.New, sha512.New} {
		key := x.Must(pls.GenerateRandomBytes(32))
		signer := x.New[*HMACSigner](WithKey(key), WithAlgorithm(hash))

		mac := hmac.New(hash, key)
		mac.Write(data)
		expectedSignature := mac.Sum(nil)

		t.Run("Sign", func(t *testing.T) {
			result := x.Must(signer.Sign(data))

			assert.NotEmpty(t, result)
			assert.Equal(t, expectedSignature, result)
		})

		t.Run("Verify", func(t *testing.T) {
			assert.True(t, signer.Verify(data, expectedSignature))

			assert.False(t, signer.Verify(data, []byte{}))
			assert.False(t, signer.Verify(data, x.Must(pls.GenerateRandomBytes(32))))
			assert.False(t, signer.Verify(x.Must(pls.GenerateRandomBytes(32)), expectedSignature))
		})
	}
}
