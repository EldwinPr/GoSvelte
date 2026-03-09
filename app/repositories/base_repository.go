package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
)

// GenerateULID creates a new unique sortable string ID.
func GenerateULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

// Pagination holds paging parameters and results.
// ... (rest of Pagination remains same)
type Pagination struct {
	Limit      int         `json:"limit"`
	Page       int         `json:"page"`
	Sort       string      `json:"sort"` // e.g., "id desc, name asc"
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		return "id DESC"
	}
	return p.Sort
}

// BaseRepository provides common CRUD operations using Go Generics.
type BaseRepository[T any] struct {
	DB *gorm.DB
}

// FilterScope for common filtering patterns
func (r *BaseRepository[T]) FilterScope(column string, operator string, value interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s %s ?", column, operator), value)
	}
}

// SearchScope for simple text searching across multiple columns
func (r *BaseRepository[T]) SearchScope(query string, columns ...string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if query == "" || len(columns) == 0 {
			return db
		}
		var conditions []string
		var values []interface{}
		for _, col := range columns {
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
			values = append(values, "%"+query+"%")
		}
		return db.Where(strings.Join(conditions, " OR "), values...)
	}
}

func (r *BaseRepository[T]) FindAll(scopes ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	var items []T
	err := r.DB.Scopes(scopes...).Find(&items).Error
	return items, err
}

func (r *BaseRepository[T]) PaginatedFind(pagination *Pagination, scopes ...func(*gorm.DB) *gorm.DB) (*Pagination, error) {
	var items []T
	var totalRows int64

	// Count total rows with filters applied
	countDb := r.DB.Model(new(T)).Scopes(scopes...)
	if err := countDb.Count(&totalRows).Error; err != nil {
		return nil, err
	}

	pagination.TotalRows = totalRows
	pagination.TotalPages = int((totalRows + int64(pagination.GetLimit()) - 1) / int64(pagination.GetLimit()))

	// Fetch data
	err := r.DB.Scopes(scopes...).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order(pagination.GetSort()).
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	pagination.Rows = items
	return pagination, nil
}

func (r *BaseRepository[T]) FindByID(id string) (*T, error) {
	var item T
	err := r.DB.First(&item, "id = ?", id).Error
	return &item, err
}

func (r *BaseRepository[T]) Create(item *T) error {
	// We use reflection/type-assertion helper logic here or handle it in service.
	// For simplicity in a generic repo, GORM hooks (BeforeCreate) are better 
	// but we'll manually ensure ID if T has an ID field.
	return r.DB.Create(item).Error
}

func (r *BaseRepository[T]) Update(item *T) error {
	return r.DB.Save(item).Error
}

func (r *BaseRepository[T]) Delete(item *T) error {
	return r.DB.Delete(item).Error
}
