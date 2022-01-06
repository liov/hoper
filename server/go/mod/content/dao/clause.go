package dao

import (
	"github.com/actliboy/hoper/server/go/mod/protobuf/content"
	"gorm.io/gorm"
)

func (c *ChainDao) MomentList(db *gorm.DB) ([]*content.Moment, error) {
	var moments []*content.Moment
	return moments, db.Scopes(c.ByName("a").ById(1).Build).Find(&moments).Error
}

func (c *ChainDao) MomentInfo(db *gorm.DB) (*content.Moment, error) {
	var moment content.Moment
	return &moment, db.Scopes(c.ByName("a").ById(1).Build).First(&moment).Error
}

func (c *ChainDao) DiaryList(db *gorm.DB) ([]*content.Diary, error) {
	var diaries []*content.Diary
	return diaries, db.Scopes(c.ByName("a").ById(1).Build).Find(&diaries).Error
}

func (c *ChainDao) DiaryInfo(db *gorm.DB) (*content.Diary, error) {
	var diaries content.Diary
	return &diaries, db.Scopes(c.ByName("a").ById(1).Build).First(&diaries).Error
}
