package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Model (One-to-One dengan Company)
type User struct {
	gorm.Model
	ID      uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Name    string    `json:"name"`
	Email   string    `json:"email" gorm:"unique"`
	Telp    string    `json:"telp"`
	Company *Company  `json:"company,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Hook BeforeCreate untuk generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return nil
}
