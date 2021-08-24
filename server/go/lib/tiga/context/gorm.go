package contexti

import (
	contexti "github.com/liov/hoper/server/go/lib/utils/context"
	"gorm.io/gorm"
)

func (c *Ctx) NewDB(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{Context: contexti.SetTranceId(c.TraceID), NewDB: true})
}
