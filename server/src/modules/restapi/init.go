package restapi

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	defaultMessage = "GhostRunner Server, HTTP REST API. Version: 0.0.1."
)

func rootEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Println("ROOT HIT") //Comment out later, for debugging purposes
	json.NewEncoder(w).Encode(infoResponse{
		Status:  http.StatusOK,
		Message: defaultMessage,
	})
}

func _initApiServer(secureServer bool, apiKey, apiCert, apiPort string) {
	apiRouter := mux.NewRouter().StrictSlash(true) // Initialize the HTTP REST API Router.

	apiRouter.HandleFunc("/", rootEndpointHandler).Methods("GET")

	if secureServer { // If a secured server is wanted. Use the specified certificate files.
		httpServer := &http.Server{
			Addr:    apiPort,   // Specify the desired HTTPS port.
			Handler: apiRouter, // Specify the above created handler.
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{ // Load the certificate and private key.
					loadTLSCertificate(apiCert, apiKey),
				},
			},
		}
		go httpServer.ListenAndServeTLS("", "")
	} else {
		go http.ListenAndServe(":"+apiPort, apiRouter) // Transform string slightly to make the expected format.
	}
}
