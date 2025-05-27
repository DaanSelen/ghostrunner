package utilities

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

const (
	InfoTag = "[INFO]"
	WarnTag = "[WARN]"
	ErrTag  = "[ERROR]"
)

func (t TokenCreateBody) GetAuthToken() string { return t.AuthToken }
func (t TokenCreateBody) GetName() string      { return t.Details.Name }

func (t TaskCreateBody) GetAuthToken() string { return t.AuthToken }
func (t TaskCreateBody) GetName() string      { return t.Details.Name }

func CheckDatabaseRemnants(databaseDir, fullDatabasePath string) {
	remnantDir := StatPath(databaseDir)
	if !remnantDir {
		log.Println(InfoTag, "Creating database folder...")
		os.Mkdir(databaseDir, os.FileMode(0755))
	}
	remnantFile := StatPath(fullDatabasePath)
	if !remnantFile {
		log.Println(InfoTag, "Creating database file...")
		os.Create(fullDatabasePath)
	}
}

func GenRandString(size int) string {
	randBytes := make([]byte, size)
	_, err := rand.Read(randBytes)
	if err != nil {
		fmt.Println(err)
	}

	randString := base64.URLEncoding.EncodeToString(randBytes)
	return randString
}

func StatPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func LoadHMACKey(keyfile string) ([]byte, error) {
	key, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, fmt.Errorf("failed to read HMAC key: %w", err)
	}
	return key, nil
}
