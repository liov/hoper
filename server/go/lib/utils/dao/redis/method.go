package redisi

import "github.com/google/uuid"

func Lock() string {
	value := uuid.New()
	cmd := "SETNX " + value.String() + " EXPIRE 100000"
	return cmd
}
