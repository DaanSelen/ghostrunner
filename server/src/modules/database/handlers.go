package database

import (
	"database/sql"
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
		return []string{}
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
		return []string{}
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

func InsertTask(name, command string, nodeids []string, date string) error {
	for _, singleNodeid := range nodeids {
		_, err := db.Exec(declStat.CreateTask, name, command, singleNodeid, date)
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}
	}
	return nil
}

func RemoveTask(name, nodeid string) error {
	var count int
	err := db.QueryRow(declStat.CountTasks, name).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count the task occurence: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("task '%s' not found", name)
	}

	if _, err = db.Exec(declStat.DeleteTask, name, nodeid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("token not found")
		}
		return fmt.Errorf("failed to retrieve token ID: %w", err)
	}

	return nil
}

func RetrieveTasks() []utilities.InternalQTaskData {
	rows, err := db.Query(declStat.ListAllTasks)
	if err != nil {
		log.Println("Query error:", err)
		return []utilities.InternalQTaskData{}
	}
	defer rows.Close()

	var tasks []utilities.InternalQTaskData

	for rows.Next() {
		var task utilities.InternalQTaskData

		err := rows.Scan(&task.Name, &task.Command, &task.Nodeid, &task.Creation)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
	}

	return tasks
}
