package db

import (
	"github.com/liov/hoper/server/go/protobuf/content"
)

func (c *ChainDao) MomentList() ([]*content.Moment, error) {
	var moments []*content.Moment
	return moments, c.db.Scopes(c.ByName("a").ById(1).Build).Find(&moments).Error
}

func (c *ChainDao) MomentInfo() (*content.Moment, error) {
	var moment content.Moment
	return &moment, c.db.Scopes(c.ByName("a").ById(1).Build).First(&moment).Error
}

func (c *ChainDao) DiaryList() ([]*content.Diary, error) {
	var diaries []*content.Diary
	return diaries, c.db.Scopes(c.ByName("a").ById(1).Build).Find(&diaries).Error
}

func (c *ChainDao) DiaryInfo() (*content.Diary, error) {
	var diaries content.Diary
	return &diaries, c.db.Scopes(c.ByName("a").ById(1).Build).First(&diaries).Error
}
