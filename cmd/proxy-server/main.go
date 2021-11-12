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
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"

	"github.com/elazarl/goproxy"
)

var (
	certificate = *flag.String("certificate", "", "Path to x509 Certificate fille")
	key         = *flag.String("key", "", "Path to the private key file")
	host        = *flag.String("host", "127.0.0.1", "Interface to bind to")
	port        = *flag.String("port", "8080", "Port to bind to")
	verbose     = *flag.Bool("verbose", false, "Enable verbose output")
)

type CertificateStore struct {
	certs map[string]*tls.Certificate
	locks map[string]*sync.Mutex
	sync.Mutex
}

func (s *CertificateStore) Fetch(host string, generate func() (*tls.Certificate, error)) (*tls.Certificate, error) {
	hostLock := s.LockFor(host)
	hostLock.Lock()
	defer hostLock.Unlock()

	cert, ok := s.certs[host]
	var err error
	if !ok {
		cert, err = generate()
		if err != nil {
			return nil, err
		}
		s.certs[host] = cert
	}
	return cert, nil
}

func (s *CertificateStore) LockFor(host string) *sync.Mutex {
	s.Lock()
	defer s.Unlock()

	lock, ok := s.locks[host]
	if !ok {
		lock = &sync.Mutex{}
		s.locks[host] = lock
	}
	return lock
}

func listenAddress() string {
	return fmt.Sprintf("%s:%s", host, port)
}

func generateSelfSignedCert() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	crtPem := &bytes.Buffer{}
	pem.Encode(crtPem, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyPem := &bytes.Buffer{}
	pem.Encode(keyPem, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return tls.X509KeyPair(crtPem.Bytes(), keyPem.Bytes())
}

func certFrom(certificate, key string) (tls.Certificate, error) {
	if certificate != "" && key != "" {
		return tls.LoadX509KeyPair(certificate, key)
	}
	return generateSelfSignedCert()
}

func main() {
	flag.Parse()

	ca, err := certFrom(certificate, key)
	if err != nil {
		log.Fatal(err)
	}
	goproxy.GoproxyCa = ca
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&ca)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&ca)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&ca)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&ca)}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	dialer := &net.Dialer{Control: func(network, address string, conn syscall.RawConn) error { return nil }}
	proxy.Tr = &http.Transport{
		Dial:            dialer.Dial,
		DialContext:     dialer.DialContext,
		TLSClientConfig: &tls.Config{},
		Proxy:           http.ProxyFromEnvironment,
	}
	proxy.CertStore = &CertificateStore{
		certs: map[string]*tls.Certificate{},
		locks: map[string]*sync.Mutex{},
	}
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func(r *http.Request, p *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		log.Printf("%s %s\n", r.Method, r.URL)
		if proxy.Verbose {
			for k, v := range r.Header {
				log.Printf("%s: %v\n", k, v)
			}
		}

		return r, nil
	})
	proxy.OnResponse().DoFunc(func(r *http.Response, p *goproxy.ProxyCtx) *http.Response {
		if r == nil {
			log.Printf("No response from server\n")
			return r
		}

		log.Printf("%d %s\n", r.StatusCode, r.Request.URL)
		if proxy.Verbose {
			for k, v := range r.Header {
				log.Printf("%s: %v\n", k, v)
			}
		}

		return r
	})

	address := listenAddress()
	log.Printf("Listening and serving HTTP on http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, proxy))
}
