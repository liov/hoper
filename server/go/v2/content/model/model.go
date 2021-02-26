package model

import "github.com/liov/hoper/go/v2/protobuf/content"

type TinyTag struct {
	Id uint64
	Name string
}

type ContentTag struct {
	RefId      uint64      `json:"refId" gorm:"size:20;not null;index:idx_content" validate:"required" annotation:"相关id"`
	Type       content.ContentType `json:"type" gorm:"type:int2;not null" validate:"required" annotation:"相关类型"`
	TagId      uint64      `json:"tagId" gorm:"size:20;not null;index:idx_content;index:idx_tag" validate:"required" annotation:"相关id"`
}