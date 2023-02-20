package main

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
	K         time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type ModelA struct {
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

const ModelTable = "model"
