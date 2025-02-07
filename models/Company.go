package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model Company (One-to-One dengan User)
type Company struct {
	ID     uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Code   string    `json:"code"`
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"user_id" gorm:"type:char(36);unique"`

	// Relasi ke User (tanpa reference langsung untuk menghindari recursive)
	User *User `json:"user,omitempty"`
}

// Hook BeforeCreate untuk generate UUID
func (c *Company) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}
