package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/xlgmokha/x/pkg/xlog"
)

func headersFrom(header http.Header) slog.Attr {
	attrs := make([]any, 0, len(header))
	for name, values := range header {
		attrs = append(attrs, slog.Any(name, values))
	}
	return slog.Group("headers", attrs...)
}

type config struct {
	certificate string
	key         string
	host        string
	port        string
	verbose     bool
}

func parseFlags(name string, args []string) (*config, error) {
	config := &config{}
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.StringVar(&config.certificate, "certificate", "", "Path to x509 certificate file")
	flags.StringVar(&config.key, "key", "", "Path to the private key file")
	flags.StringVar(&config.host, "host", "127.0.0.1", "Interface to bind to")
	flags.StringVar(&config.port, "port", "8080", "Port to bind to")
	flags.BoolVar(&config.verbose, "verbose", false, "Enable verbose output")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *config) listenAddress() string {
	return net.JoinHostPort(c.host, c.port)
}

type CertificateStore struct {
	mu    sync.Mutex
	certs map[string]*tls.Certificate
	locks map[string]*sync.Mutex
}

func newCertificateStore() *CertificateStore {
	return &CertificateStore{
		certs: map[string]*tls.Certificate{},
		locks: map[string]*sync.Mutex{},
	}
}

func (s *CertificateStore) Fetch(hostname string, gen func() (*tls.Certificate, error)) (*tls.Certificate, error) {
	hostLock := s.lockFor(hostname)
	hostLock.Lock()
	defer hostLock.Unlock()

	if cert, ok := s.get(hostname); ok {
		return cert, nil
	}

	cert, err := gen()
	if err != nil {
		return nil, err
	}
	s.put(hostname, cert)
	return cert, nil
}

func (s *CertificateStore) get(hostname string) (*tls.Certificate, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cert, ok := s.certs[hostname]
	return cert, ok
}

func (s *CertificateStore) put(hostname string, cert *tls.Certificate) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.certs[hostname] = cert
}

func (s *CertificateStore) lockFor(hostname string) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()

	lock, ok := s.locks[hostname]
	if !ok {
		lock = &sync.Mutex{}
		s.locks[hostname] = lock
	}
	return lock
}

func generateSelfSignedCert(host string) (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"xlg"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	crtPem := &bytes.Buffer{}
	pem.Encode(crtPem, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyPem := &bytes.Buffer{}
	pem.Encode(keyPem, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return tls.X509KeyPair(crtPem.Bytes(), keyPem.Bytes())
}

func certFrom(certificate, key, host string) (tls.Certificate, error) {
	if certificate != "" && key != "" {
		return tls.LoadX509KeyPair(certificate, key)
	}
	return generateSelfSignedCert(host)
}

func main() {
	logger := xlog.New(os.Stdout, xlog.Fields{})

	config, err := parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		logger.Error("invalid flags", slog.Any("error", err))
		os.Exit(1)
	}

	ca, err := certFrom(config.certificate, config.key, config.host)
	if err != nil {
		logger.Error("could not load a certificate", slog.Any("error", err))
		os.Exit(1)
	}
	goproxy.GoproxyCa = ca
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&ca)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&ca)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&ca)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&ca)}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = config.verbose
	dialer := &net.Dialer{Control: func(network, address string, conn syscall.RawConn) error { return nil }}
	proxy.Tr = &http.Transport{
		Dial:            dialer.Dial,
		DialContext:     dialer.DialContext,
		TLSClientConfig: &tls.Config{},
		Proxy:           http.ProxyFromEnvironment,
	}
	proxy.CertStore = newCertificateStore()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func(r *http.Request, p *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
		}
		if proxy.Verbose {
			attrs = append(attrs, headersFrom(r.Header))
		}
		logger.LogAttrs(r.Context(), slog.LevelInfo, "request", attrs...)

		return r, nil
	})
	proxy.OnResponse().DoFunc(func(r *http.Response, p *goproxy.ProxyCtx) *http.Response {
		if r == nil {
			logger.Warn("no response from server")
			return r
		}

		attrs := []slog.Attr{
			slog.Int("status", r.StatusCode),
			slog.String("url", r.Request.URL.String()),
		}
		if proxy.Verbose {
			attrs = append(attrs, headersFrom(r.Header))
		}
		logger.LogAttrs(r.Request.Context(), slog.LevelInfo, "response", attrs...)

		return r
	})

	address := config.listenAddress()
	logger.Info("listening", slog.String("address", address))

	if err := http.ListenAndServe(address, proxy); err != nil {
		logger.Error("server stopped", slog.Any("error", err))
		os.Exit(1)
	}
}
