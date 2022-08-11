package clausei

import "gorm.io/gorm"

func First[T any](db *gorm.DB) (*T, error) {
	t := new(T)
	err := db.First(t).Error
	return t, err
}

type DB[T any] gorm.DB

func (db *DB[T]) First() (*T, error) {
	t := new(T)
	err := (*gorm.DB)(db).First(t).Error
	return t, err
}
