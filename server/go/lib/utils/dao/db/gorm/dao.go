package gormi

import (
	"context"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	ChainDao
}

func NewRepository[T any](ctx context.Context, db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		*NewChainDao(ctx, db),
	}
}

func (r *Repository[T]) Create(t *T) error {
	return r.DB.Create(t).Error
}

func (r *Repository[T]) Read(id int) (*T, error) {
	var t T
	err := r.DB.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository[T]) Update(t *T) error {
	return r.DB.Updates(&t).Error
}

func (r *Repository[T]) Delete(id int) error {
	var t T
	return r.DB.Delete(&t, id).Error
}
