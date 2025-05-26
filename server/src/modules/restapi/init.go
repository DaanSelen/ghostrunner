package restapi

import (
	"encoding/json"
	"errors"
	"ghostrunner-server/modules/utilities"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	defaultMessage   = "GhostRunner Server, HTTP REST API. Version: 0.0.1."
	readWriteTimeout = 30 * time.Second //Seconds
)

func rootEndpointHandler(w http.ResponseWriter, r *http.Request) { // This endpoint handles has been placed in the init section because its basic.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Println("Root HTTP API endpoint has been reached.") //Comment out later, for debugging purposes
	json.NewEncoder(w).Encode(utilities.InfoResponse{
		Status:  http.StatusOK,
		Message: defaultMessage,
	})
}

func InitApiServer(cfg utilities.ConfigStruct, hmacKey []byte) {
	rtr := createRouter(hmacKey)
	srv := createServer(cfg, rtr)

	// Following func can be goroutines.
	go func() {
		var err error
		if cfg.Secure {
			if utilities.StatPath(cfg.ApiCertFile) && utilities.StatPath(cfg.ApiKeyFile) {
				err = srv.ListenAndServeTLS(cfg.ApiCertFile, cfg.ApiKeyFile)
			} else {
				err = errors.New("failed to find one or both certificate- and/or keyfile")
			}
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil {
			log.Println(utilities.ErrTag, err)
		}
		defer srv.Close()
	}()
	log.Println(utilities.InfoTag, "Successfully started the GhostServer goroutine at:", cfg.Address)
}

func createRouter(hmacKey []byte) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", rootEndpointHandler).Methods("GET")

	r.HandleFunc("/token/create", createTokenHandler(hmacKey)).Methods("POST")
	r.HandleFunc("/token/delete", deleteTokenHandler(hmacKey)).Methods("DELETE")
	r.HandleFunc("/token/list", listTokenHandler(hmacKey)).Methods("GET")

	r.HandleFunc("/task/create", createTaskHandler(hmacKey)).Methods("POST")
	r.HandleFunc("/task/delete", deleteTaskHandler(hmacKey)).Methods("DELETE")
	r.HandleFunc("/task/list", listTasksHandler(hmacKey)).Methods("GET")

	return r
}

func createServer(cfg utilities.ConfigStruct, ghostHandler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.Address,  // Specify the desired HTTPS port.
		Handler:      ghostHandler, // Specify the above created handler.
		ReadTimeout:  readWriteTimeout,
		WriteTimeout: readWriteTimeout,
		ErrorLog:     log.New(io.Discard, "", 0),
	}
}
