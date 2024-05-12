package db

import (
	"github.com/liov/hoper/server/go/protobuf/content"
)

func (c *ChainDao) MomentList() ([]*content.Moment, error) {
	var moments []*content.Moment
	return moments, c.ByName("a").ById(1).Exec(c.db).Find(&moments).Error
}

func (c *ChainDao) MomentInfo() (*content.Moment, error) {
	var moment content.Moment
	return &moment, c.ByName("a").ById(1).Exec(c.db).First(&moment).Error
}

func (c *ChainDao) DiaryList() ([]*content.Diary, error) {
	var diaries []*content.Diary
	return diaries, c.ByName("a").ById(1).Exec(c.db).Find(&diaries).Error
}

func (c *ChainDao) DiaryInfo() (*content.Diary, error) {
	var diary content.Diary
	return &diary, c.ByName("a").ById(1).Exec(c.db).First(&diary).Error
}
