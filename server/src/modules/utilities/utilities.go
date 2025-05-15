package utilities

import (
	"crypto/tls"
	"log"
)

func HandleError(err error, task string) {
	if err != nil {
		log.Fatal("The program crashed unexpectedly while doing: "+task+"\nThe following exception occured:", err)
	}
}

func ConsoleLog(message string) {
	log.Println(message)
}

func LoadCertificate(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	return cert
}
