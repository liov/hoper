package model

import "time"

type User struct {
	ID          uint64     `gorm:"primary_key" json:"id"`
	ActivatedAt *time.Time `json:"-"` //激活时间
	Name        string     `gorm:"type:varchar(10);not null" json:"name"`
	Password    string     `gorm:"type:varchar(100)" json:"-"`
	Account     string     `gorm:"type:varchar(20);unique_index" json:"-"`
	Email       string     `gorm:"type:varchar(20);unique_index;not null" json:"email"`
	Phone       *string    `gorm:"type:varchar(20);unique_index" json:"phone"` //手机号
	CreatedAt   time.Time  `json:"-"`
	Status      uint8      `gorm:"type:smallint;default:0" json:"-"` //状态
}
