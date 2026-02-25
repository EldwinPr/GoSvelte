package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primaryKey;type:char(26)" json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Clearance int    `json:"clearance"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = GenerateULID()
	return nil
}
