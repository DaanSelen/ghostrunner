package database

import (
	"database/sql"
	"errors"
	"ghostrunner-server/modules/utilities"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseDir      = "./data"
	fullDatabasePath = databaseDir + "/ghostserver.db"
)

var db *sql.DB

func InitSqlite(adminToken string, hmacKey []byte) {
	utilities.CheckDatabaseRemnants(databaseDir, fullDatabasePath)

	var err error
	db, err = sql.Open("sqlite3", fullDatabasePath)
	if err != nil {
		log.Println(utilities.ErrTag, err)
	}

	_, err = db.Exec(declStat.SetupDatabase)
	if err != nil {
		log.Println(utilities.ErrTag, err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(utilities.ErrTag, err)
	}

	var adminTokenID int
	err = db.QueryRow("SELECT id FROM tokens WHERE id = '0'").Scan(&adminTokenID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(utilities.InfoTag, "No Admin token detected, inserting...")
			insertAdminToken(adminToken, hmacKey)
		} else {
			log.Println(utilities.InfoTag, "Something else went wrong, doing nothing.")
		}
	} else {
		log.Println(utilities.InfoTag, "An Admin token is already present, not re-inserting.")
	}
}
