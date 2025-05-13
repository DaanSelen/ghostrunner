package restapi

import (
	"crypto/tls"
	"log"
)

func loadTLSCertificate(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	return cert
}
