package dbi

import "time"

type TimeModel struct {
	CreatedAt time.Time  `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at"  gorm:"default:current_timestamp"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type TimeStampModel struct {
	CreatedAt int64 `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt int64 `json:"updated_at" gorm:"default:current_timestamp"`
	DeletedAt int64 `json:"deleted_at"`
}

type TimeStringModel struct {
	CreatedAt string `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt string `json:"updated_at" gorm:"default:current_timestamp"`
	DeletedAt string `json:"deleted_at"`
}
