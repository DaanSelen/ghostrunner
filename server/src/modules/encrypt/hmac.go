package encrypt

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"errors"
)

func CreateHMAC(token string, key []byte) string {
	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(token))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func ValidateHMAC(expectedHMAC, candToken string, key []byte) (bool, error) {
	if len(candToken) == 0 {
		return false, errors.New("candidate MAC is empty")
	}
	candMac := CreateHMAC(candToken, key)

	return hmac.Equal([]byte(expectedHMAC), []byte(candMac)), nil
}
