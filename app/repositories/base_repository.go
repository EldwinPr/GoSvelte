package repositories

import "gorm.io/gorm"

// BaseRepository provides common CRUD operations using Go Generics.
type BaseRepository[T any] struct {
	DB *gorm.DB
}

func (r *BaseRepository[T]) FindAll() ([]T, error) {
	var items []T
	err := r.DB.Find(&items).Error
	return items, err
}

func (r *BaseRepository[T]) FindByID(id uint) (*T, error) {
	var item T
	err := r.DB.First(&item, id).Error
	return &item, err
}

func (r *BaseRepository[T]) Create(item *T) error {
	return r.DB.Create(item).Error
}

func (r *BaseRepository[T]) Update(item *T) error {
	return r.DB.Save(item).Error
}

func (r *BaseRepository[T]) Delete(item *T) error {
	return r.DB.Delete(item).Error
}
