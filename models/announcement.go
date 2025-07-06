package models

import "gorm.io/gorm"

// Announcement represents the announcement table
type Announcement struct {
	gorm.Model
	Name         string  `json:"name" gorm:"not null;index" binding:"required"`
	Description  *string `json:"description" gorm:"type:text"`
	ExternalLink *string `json:"external_link" gorm:"type:text"`
	HasRead      bool    `json:"has_read" gorm:"default:false;index"`
}
