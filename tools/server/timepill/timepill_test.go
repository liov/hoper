package timepill

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	var comments []*Comment
	err := json.Unmarshal([]byte(`[]`), &comments)
	fmt.Println(err, comments)
}

func TestGetDiaryComments(t *testing.T) {
	Token = "Basic bGJ5LmlAcXEuY29tOmxieTYwNA=="
	fmt.Println(ApiService.GetDiaryComments(6817247))
}
