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

func (r *BaseRepository[T]) FindByID(id any) (*T, error) {
	var item T
	err := r.DB.First(&item, "id = ?", id).Error
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

type PaginationResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
}

func (r *BaseRepository[T]) Paginate(page, pageSize int, query *gorm.DB) (*PaginationResult[T], error) {
	var items []T
	var totalCount int64

	if query == nil {
		query = r.DB
	}

	// Get total count before pagination
	if err := query.Model(new(T)).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return &PaginationResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
