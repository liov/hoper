package clause

import (
	"context"
	"gorm.io/gorm"
)

type TestChainDBDao struct {
	ChainDB
}

func NewTestChainDBDao(ctx context.Context) *TestChainDBDao {
	return &TestChainDBDao{NewChainDB(ctx)}
}

func testChainDBDao() {
	dao := NewTestChainDBDao(context.Background())
	db := new(gorm.DB)
	dao.ById(1).ByName("a").List(db)
}
