package timepill

import (
	"context"
	"fmt"
	"github.com/liov/hoper/server/go/lib/initialize"
	_type "github.com/liov/hoper/server/go/lib/utils/dao/db/type"
	"testing"
)

func TestTable(t *testing.T) {
	defer initialize.Start(&Conf, &Dao)()
	dao := &DBDao{ctx: context.Background(), Hoper: Dao.Hoper.DB}
	diaries, err := dao.List(_type.NewListReq[int](1, 10))
	if err != nil {
		t.Error(err)
	}
	for _, diary := range diaries {
		fmt.Println(diary.Content)
	}
}
