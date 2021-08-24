package model

import (
	"github.com/liov/hoper/server/go/mod/protobuf/content"
)

type TinyTag struct {
	Id   uint64
	Name string
}

type ContentTag struct {
	RefId uint64              `json:"refId" gorm:"size:20;not null;index:idx_content" validate:"required" annotation:"相关id"`
	Type  content.ContentType `json:"type" gorm:"type:int2;not null" validate:"required" annotation:"相关类型"`
	TagId uint64              `json:"tagId" gorm:"size:20;not null;index:idx_content;index:idx_tag" validate:"required" annotation:"相关id"`
}

type ContentTagRel struct {
	RefId uint64 `json:"refId" validate:"required" annotation:"相关id"`
	content.TinyTag
}

type ContentExt struct {
	Type  content.ContentType
	RefId uint64
}

type ContentAction struct {
	Id     uint64
	Type   content.ContentType
	RefId  uint64 `json:"refId" validate:"required" annotation:"相关id"`
	Action content.ActionType
}

type ContentCollect struct {
	Id    uint64
	Type  content.ContentType
	RefId uint64 `json:"refId" validate:"required" annotation:"相关id"`
	FavId uint64 `json:"favId" validate:"required" annotation:"收藏夹id"`
}
type Collect struct {
	Id     uint64
	Type   content.ContentType
	RefId  uint64
	UserId uint64
	FavId  uint64
}
