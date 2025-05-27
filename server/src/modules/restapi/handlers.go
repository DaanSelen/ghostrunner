package restapi

import (
	"encoding/json"
	"ghostrunner-server/modules/database"
	"ghostrunner-server/modules/encrypt"
	"ghostrunner-server/modules/utilities"
	"log"
	"net/http"
	"strings"
	"time"

	"slices"
)

type authPayload interface {
	GetAuthToken() string
	GetName() string
}

func generalAuth(w http.ResponseWriter, securedCandidate string) bool {
	tokens := database.RetrieveTokens()
	if !slices.Contains(tokens, securedCandidate) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}

func parseAndAuth[T authPayload](w http.ResponseWriter, r *http.Request, hmacKey []byte) (T, bool) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(utilities.ErrTag, "Decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return data, false
	}

	if data.GetAuthToken() == "" || data.GetName() == "" {
		log.Println("[ERROR] Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return data, false
	}

	securedCandidate := encrypt.CreateHMAC(data.GetAuthToken(), hmacKey)
	return data, generalAuth(w, securedCandidate)
}

/*
The following section pertrains to Token creation and deletion.
*/

func createTokenHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := parseAndAuth[utilities.TokenCreateBody](w, r, hmacKey)
		if !ok {
			return
		}

		data.Details.Name = strings.ToLower(data.Details.Name) //Transform to lower
		token, err := createToken(data.Details.Name, hmacKey)
		if err != nil {
			log.Println(utilities.ErrTag, "createToken failed:", err)
			http.Error(w, "Token creation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(utilities.InfoResponse{
			Status:  http.StatusCreated,
			Message: "Token Successfully Created.",
			Data:    token,
		})
	}
}

func deleteTokenHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := parseAndAuth[utilities.TokenCreateBody](w, r, hmacKey)
		if !ok {
			return
		}

		if err := deleteToken(data.Details.Name, hmacKey); err != nil {
			log.Println(utilities.ErrTag, "deleteToken failed:", err)
			http.Error(w, "Token deletion failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utilities.InfoResponse{
			Status:  http.StatusOK,
			Message: "Token Deleted Successfully",
		})
	}
}

func listTokenHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenCandidate := authHeader[len(prefix):]
		securedCandidate := encrypt.CreateHMAC(tokenCandidate, hmacKey)
		if !generalAuth(w, securedCandidate) {
			return
		}

		data := database.RetrieveTokenNames()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utilities.InfoResponse{
			Status:  http.StatusOK,
			Message: "Successfully Retrieved Tokens",
			Data:    data,
		})
	}
}

func createToken(tokenName string, hmacKey []byte) (string, error) {
	randomString := utilities.GenRandString(64)
	securedString := encrypt.CreateHMAC(randomString, hmacKey)
	tokenName = strings.ToLower(tokenName)

	if err := database.InsertToken(tokenName, securedString); err != nil {
		return "", err
	}
	return randomString, nil
}

func deleteToken(tokenName string, hmacKey []byte) error {
	return database.RemoveToken(tokenName)
}

/*
The following section portrains to Task creation and deletion.
*/

func createTaskHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := parseAndAuth[utilities.TaskCreateBody](w, r, hmacKey)
		if !ok {
			return
		}

		data.Details.Name = strings.ToLower(data.Details.Name) //Transform to lower
		if err := createTask(data.Details.Name, data.Details.Command, data.Details.Nodeids); err != nil {
			log.Println(utilities.ErrTag, "createTask failed:", err)
			http.Error(w, "Task creation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utilities.InfoResponse{
			Status:  http.StatusOK,
			Message: "Task '" + data.Details.Name + "' Created Successfully.",
		})
	}
}

func deleteTaskHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := parseAndAuth[utilities.TaskCreateBody](w, r, hmacKey)
		if !ok {
			return
		}
		nodeid := data.Details.Nodeids[0]

		if err := deleteTask(data.Details.Name, nodeid); err != nil {
			log.Println(utilities.ErrTag, "createTask failed:", err)
			http.Error(w, "Task deletion failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utilities.InfoResponse{
			Status:  http.StatusOK,
			Message: "Task '" + data.Details.Name + "' Deleted Successfully.",
		})
	}
}

func listTasksHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenCandidate := authHeader[len(prefix):]
		securedCandidate := encrypt.CreateHMAC(tokenCandidate, hmacKey)
		if !generalAuth(w, securedCandidate) {
			return
		}

		data := database.RetrieveTasks()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utilities.InfoResponse{
			Status:  http.StatusOK,
			Message: "Successfully Retrieved Tasks",
			Data:    data,
		})
	}
}

func flushTaskListHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("skibidi")
	}
}

func createTask(taskName, command string, nodeids []string) error {
	creationDate := time.Now().Format("02-01-2006 15:04:05")
	taskName = strings.ToLower(taskName)

	return database.InsertTask(taskName, command, nodeids, creationDate)
}

func deleteTask(taskName, nodeid string) error {
	return database.RemoveTask(taskName, nodeid)
}
