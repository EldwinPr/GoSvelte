package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primaryKey;type:char(26)" json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Clearance int    `json:"clearance"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = GenerateULID()
	return nil
}
