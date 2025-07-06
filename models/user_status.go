package models

import "gorm.io/gorm"

// UserStatus represents the status table
type UserStatus struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null;index" binding:"required"`

	// Reverse relationships
	Users []User `json:"users,omitempty" gorm:"foreignKey:StatusID"`
}
