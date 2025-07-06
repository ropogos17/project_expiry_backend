package models

import (
	"gorm.io/gorm"
)

// Role represents the role table
type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null;index" binding:"required"`

	// Reverse relationships
	Users []User `json:"users,omitempty" gorm:"foreignKey:RoleID"`
}
