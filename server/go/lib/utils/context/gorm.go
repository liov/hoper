package contexti

import (
	"gorm.io/gorm"
)

func (c *RequestContext) NewDB(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{Context: SetTranceId(c.TraceID), NewDB: true})
}
