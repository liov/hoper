package model

import "gorm.io/gorm"

type User struct {
	ID        int64
	Name      string
	Mail      string
	Phone     string
	DeletedAt gorm.DeletedAt
}
