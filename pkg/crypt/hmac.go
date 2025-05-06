package crypt

import (
	"crypto/hmac"
	"hash"

	"github.com/xlgmokha/x/pkg/x"
)

type HMACSigner struct {
	key     []byte
	factory func() hash.Hash
}

func (s *HMACSigner) Sign(data []byte) ([]byte, error) {
	mac := hmac.New(s.factory, s.key)
	_, err := mac.Write(data)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
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
