package gorm

import "gorm.io/gorm"

func GetById[T any](db *gorm.DB, id int) (*T, error) {
	t := new(T)
	err := db.First(t, id).Error
	return t, err
}

type DB[T any] gorm.DB

func (db *DB[T]) GetById(id int) (*T, error) {
	t := new(T)
	err := (*gorm.DB)(db).First(t, id).Error
	return t, err
}

type DB2[T any] struct {
	gorm.DB
}

func (db *DB2[T]) GetById(id int) (*T, error) {
	t := new(T)
	err := db.First(t, id).Error
	return t, err
}
