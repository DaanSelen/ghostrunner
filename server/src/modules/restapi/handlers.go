package restapi

import (
	"ghostrunner-server/modules/utilities"
	"net/http"

	"github.com/gorilla/mux"
)

func handleSchedule(w http.ResponseWriter, r *http.Request) {
	action := mux.Vars(r)["action"]

	utilities.ConsoleLog("Funky Funky " + action)
}
