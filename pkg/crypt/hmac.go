package crypt

import (
	"crypto/hmac"
	"hash"

	"github.com/xlgmokha/x/pkg/x"
)

type HMACSigner struct {
	key     []byte
	factory x.Factory[hash.Hash]
}

func WithAlgorithm(factory x.Factory[hash.Hash]) x.Option[*HMACSigner] {
	return x.With[*HMACSigner](func(item *HMACSigner) {
		item.factory = factory
	})
}

func WithKey(key []byte) x.Option[*HMACSigner] {
	return x.With[*HMACSigner](func(item *HMACSigner) {
		item.key = key
	})
}

func (s *HMACSigner) Sign(data []byte) ([]byte, error) {
	mac := hmac.New(s.factory, s.key)
	_, err := mac.Write(data)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func (s *HMACSigner) Verify(data []byte, signature []byte) bool {
	actual, err := s.Sign(data)
	if err != nil {
		return false
	}

	return hmac.Equal(actual, signature)
}
