package auth

import (
	"encoding/base64"
	"net/http"
)

func SetBasicAuth(header http.Header, username, password string) {
	header.Set("Authorization", "Basic "+BasicAuth(username, password))
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
