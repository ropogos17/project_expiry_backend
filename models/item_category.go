package models

import "gorm.io/gorm"

// ItemCategory represents the category table
type ItemCategory struct {
	gorm.Model
	Name string `json:"name" gorm:"not null;index" binding:"required"`

	// Reverse relationships
	Items []Item `json:"items,omitempty" gorm:"foreignKey:CategoryID"`
}
