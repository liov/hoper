package gorm

import "time"

type Model struct {
	Id        int
	A         int
	B         int
	C         int
	D         int
	E         int
	F         string
	H         string
	I         string
	J         string
	K         string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
