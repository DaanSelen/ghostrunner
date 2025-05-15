package restapi

import (
	"encoding/json"
	"ghostrunner-server/modules/confread"
	"ghostrunner-server/modules/utilities"
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

	utilities.ConsoleLog("ROOT HIT") //Comment out later, for debugging purposes
	json.NewEncoder(w).Encode(infoResponse{
		Status:  http.StatusOK,
		Message: defaultMessage,
	})
}

func InitApiServer(cfg confread.ConfigStruct) {
	rtr := createRouter()
	srv := createServer(cfg, rtr)

	go func() {
		var err error
		if cfg.Secure {
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		utilities.HandleError(err, "Initializing the HTTP REST API!")
	}()
	utilities.ConsoleLog("Successfully started the GhostServer goroutine.")
}

func createRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", rootEndpointHandler).Methods("GET")
	r.HandleFunc("/schedule/{action:register|deregister}", handleSchedule).Methods("POST")

	return r
}

func createServer(cfg confread.ConfigStruct, ghostHandler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.Address,  // Specify the desired HTTPS port.
		Handler:      ghostHandler, // Specify the above created handler.
		ReadTimeout:  readWriteTimeout,
		WriteTimeout: readWriteTimeout,
	}
}
