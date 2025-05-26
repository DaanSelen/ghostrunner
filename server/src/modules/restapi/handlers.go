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

const (
	constCreationStatus string = "Created"
)

func generalAuth(w http.ResponseWriter, securedCandidate string) bool {
	tokens := database.RetrieveTokens()
	if !slices.Contains(tokens, securedCandidate) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}

func parseTokenAndAuth(w http.ResponseWriter, r *http.Request, hmacKey []byte) (utilities.TokenCreateBody, bool) {
	var data utilities.TokenCreateBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(utilities.ErrTag, "Decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return data, false
	}

	if data.AuthToken == "" || data.Details.Name == "" {
		log.Println("[ERROR] Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return data, false
	}

	givenToken := data.AuthToken
	securedCandidate := encrypt.CreateHMAC(givenToken, hmacKey)
	return data, generalAuth(w, securedCandidate)
}

func parseTaskAndAuth(w http.ResponseWriter, r *http.Request, hmacKey []byte) (utilities.TaskBody, bool) {
	var data utilities.TaskBody
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(utilities.ErrTag, "Decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return data, false
	}

	if data.AuthToken == "" || data.Details.Name == "" {
		log.Println("[ERROR] Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return data, false
	}

	givenToken := data.AuthToken
	securedCandidate := encrypt.CreateHMAC(givenToken, hmacKey)
	return data, generalAuth(w, securedCandidate)
}

/*
The following section portrains to Token creation and deletion.
*/

func createTokenHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := parseTokenAndAuth(w, r, hmacKey)
		if !ok {
			return
		}

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
			Message: "Token Succesfully Created.",
			Data:    token,
		})
	}
}

func deleteTokenHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := parseTokenAndAuth(w, r, hmacKey)
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
			Message: "Succesfully Retrieved Tokens",
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
		data, ok := parseTaskAndAuth(w, r, hmacKey)
		if !ok {
			return
		}

		if err := createTask(data.Details.Name, data.Details.Command, data.Details.Nodeids); err != nil {
			log.Println(utilities.ErrTag, "createTask failed:", err)
			http.Error(w, "Task creation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utilities.InfoResponse{
			Status:  http.StatusOK,
			Message: "Task '" + data.Details.Name + "' Created Succesfully.",
		})
	}
}

func deleteTaskHandler(hmacKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := parseTaskAndAuth(w, r, hmacKey)
		if !ok {
			return
		}

		if err := deleteTask(data.Details.Name); err != nil {
			log.Println(utilities.ErrTag, "createTask failed:", err)
			http.Error(w, "Task deletion failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utilities.InfoResponse{
			Status:  http.StatusOK,
			Message: "Task '" + data.Details.Name + "' Deleted Succesfully.",
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
			Message: "Succesfully Retrieved Tasks",
			Data:    data,
		})
	}
}

func createTask(taskName, command string, nodeids []string) error {
	creationDate := time.Now().Format("02-01-2006 15:04:05")
	creationStatus := constCreationStatus
	taskName = strings.ToLower(taskName)

	return database.InsertTask(taskName, command, nodeids, creationDate, creationStatus)
}

func deleteTask(taskName string) error {
	return database.RemoveTask(taskName)
}
