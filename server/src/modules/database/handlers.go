package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"ghostrunner-server/modules/encrypt"
	"ghostrunner-server/modules/utilities"
	"log"
	"strings"
)

func insertAdminToken(token string, hmacKey []byte) error {
	var adminTokenName string = "Self-Generated Admin Token"
	adminTokenName = strings.ToLower(adminTokenName)
	hashedToken := encrypt.CreateHMAC(token, hmacKey)

	_, err := db.Exec(declStat.AdminTokenCreate, adminTokenName, hashedToken)
	if err != nil {
		return fmt.Errorf("failed to insert admin token: %w", err)
	}

	return nil
}

func InsertToken(tokenName, securedToken string) error {
	_, err := db.Exec(declStat.CreateToken, tokenName, securedToken)
	return err
}

func RemoveToken(tokenName string) error {
	var tokenID int
	err := db.QueryRow(declStat.GetTokenID, tokenName).Scan(&tokenID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("token not found")
		}
		return fmt.Errorf("failed to retrieve token ID: %w", err)
	}

	if tokenID == 0 { //TRYING TO REMOVE THE ADMIN TOKEN! NOT ALLOWED!
		return fmt.Errorf("not abiding the removal of the admin token, program resisted")
	}

	_, err = db.Exec(declStat.DeleteToken, tokenName)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

func RetrieveTokens() []string {
	rows, err := db.Query(declStat.RetrieveTokens)
	if err != nil {
		log.Println(utilities.ErrTag, err)
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var singleToken string
		err = rows.Scan(&singleToken)
		if err != nil {
			log.Println(utilities.ErrTag, err)
		}

		tokens = append(tokens, singleToken)
	}
	return tokens
}

func RetrieveTokenNames() []string {
	rows, err := db.Query(declStat.RetrieveTokenNames)
	if err != nil {
		log.Println(utilities.ErrTag, err)
	}
	defer rows.Close()

	var tokenNames []string
	for rows.Next() {
		var singleTokenName string
		err = rows.Scan(&singleTokenName)
		if err != nil {
			log.Println(utilities.ErrTag, err)
		}

		tokenNames = append(tokenNames, singleTokenName)
	}
	return tokenNames
}

func InsertTask(name, command string, nodeids []string, date, status string) error {
	for _, singleNodeid := range nodeids {
		_, err := db.Exec(declStat.CreateTask, name, command, string(singleNodeid), date, status)
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}
	}
	return nil
}

func RemoveTask(name string) error {
	_, err := db.Exec(declStat.DeleteTask, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("token not found")
		}
		return fmt.Errorf("failed to retrieve token ID: %w", err)
	}

	return nil
}

func RetrieveTasks() []utilities.TaskData {
	rows, err := db.Query(declStat.ListAllTasks)
	if err != nil {
		log.Println("Query error:", err)
		return nil
	}
	defer rows.Close()

	var tasks []utilities.TaskData

	for rows.Next() {
		var task utilities.TaskData
		var nodeidsStr string

		err := rows.Scan(&task.Name, &task.Command, &nodeidsStr, &task.Creation, &task.Status)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		err = json.Unmarshal([]byte(nodeidsStr), &task.Nodeids)
		if err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
	}

	return tasks
}
