package models

import "gorm.io/gorm"

// ActionType represents the action type table
type ActionType struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null;index" binding:"required"`

	// Reverse relationships
	Logs []Log `json:"logs,omitempty" gorm:"foreignKey:ActionTypeID"`
}
