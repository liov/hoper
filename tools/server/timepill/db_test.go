package timepill

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"testing"
)

func TestTable(t *testing.T) {
	defer initialize.Start(&Conf, &Dao)()
	CreateFaceTable()
}
