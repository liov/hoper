package clausei

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmptyClause struct {
}

func (receiver *EmptyClause) Build(builder clause.Builder) {

}

func (receiver *EmptyClause) ModifyStatement(*gorm.Statement) {

}
