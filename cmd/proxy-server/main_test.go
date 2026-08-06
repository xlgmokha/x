package main

import (
	"crypto/tls"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	t.Run("applies the supplied values", func(t *testing.T) {
		config, err := parseFlags("proxy-server", []string{"-host", "0.0.0.0", "-port", "9999", "-verbose"})

		require.NoError(t, err)
		assert.Equal(t, "0.0.0.0", config.host)
		assert.Equal(t, "9999", config.port)
		assert.True(t, config.verbose)
	})

	t.Run("falls back to the defaults", func(t *testing.T) {
		config, err := parseFlags("proxy-server", []string{})

		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1", config.host)
		assert.Equal(t, "8080", config.port)
		assert.False(t, config.verbose)
		assert.Equal(t, "", config.certificate)
		assert.Equal(t, "", config.key)
	})

	t.Run("builds the listen address", func(t *testing.T) {
		config, err := parseFlags("proxy-server", []string{"-host", "0.0.0.0", "-port", "9999"})

		require.NoError(t, err)
		assert.Equal(t, "0.0.0.0:9999", config.listenAddress())
	})

	t.Run("brackets an ipv6 literal", func(t *testing.T) {
		config, err := parseFlags("proxy-server", []string{"-host", "::1"})

		require.NoError(t, err)
		assert.Equal(t, "[::1]:8080", config.listenAddress())
	})
}

func TestCertificateStore(t *testing.T) {
	generate := func() (*tls.Certificate, error) { return &tls.Certificate{}, nil }

	t.Run("concurrent fetches for different hosts do not race", func(t *testing.T) {
		store := newCertificateStore()

		var wg sync.WaitGroup
		for _, host := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
			wg.Add(1)
			go func(host string) {
				defer wg.Done()
				_, _ = store.Fetch(host, generate)
			}(host)
		}

		wg.Wait()
	})

	t.Run("generates once per host", func(t *testing.T) {
		store := newCertificateStore()
		calls := 0
		counting := func() (*tls.Certificate, error) {
			calls++
			return &tls.Certificate{}, nil
		}

		first, err := store.Fetch("example.com", counting)
		require.NoError(t, err)

		second, err := store.Fetch("example.com", counting)
		require.NoError(t, err)

		assert.Same(t, first, second)
		assert.Equal(t, 1, calls)
	})
}
