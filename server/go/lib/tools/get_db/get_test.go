package get_db

import (
	"testing"
	"time"
)

func Test_Time(t *testing.T) {
	now := time.Now()
	GetDB().Exec(`insert into education (deleted_at) values (?)`, now)
}
