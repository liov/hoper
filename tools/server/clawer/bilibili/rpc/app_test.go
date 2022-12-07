package rpc

import (
	"log"
	"testing"
)

func TestAPP(t *testing.T) {
	log.Println(Get[*VideoInfo](GetPlayerUrlV2(571773475, 120)))
}
