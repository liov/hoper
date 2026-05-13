package icebridge

import (
	"crypto/rand"
	"encoding/hex"
)

// NewCredentials 生成 ICE ufrag/pwd。
func NewCredentials() (ufrag, pwd string, err error) {
	ufrag, err = randToken(8)
	if err != nil {
		return "", "", err
	}
	pwd, err = randToken(24)
	return ufrag, pwd, err
}

func randToken(n int) (string, error) {
	b := make([]byte, (n+1)/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := hex.EncodeToString(b)
	if len(s) < n {
		return s, nil
	}
	return s[:n], nil
}
