package models

import "gorm.io/gorm"

// ShoppingListItem represents the shopping list item table
type ShoppingListItem struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null;index" binding:"required"`
	Description *string `json:"description" gorm:"type:text"`
	Quantity    *string `json:"quantity"`
	IsChecked   bool    `json:"is_checked" gorm:"default:false;index"`
	IsDeleted   bool    `json:"is_deleted" gorm:"default:false;index"`

	// User relationship
	LastModified uint `json:"last_modified" gorm:"not null;index"`
	ModifiedBy   User `json:"modified_by" gorm:"foreignKey:LastModified;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// Scopes for common queries
func (ShoppingListItem) NotDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("is_deleted = ?", false)
}

func (ShoppingListItem) Checked(db *gorm.DB) *gorm.DB {
	return db.Where("is_checked = ?", true)
}

func (ShoppingListItem) Unchecked(db *gorm.DB) *gorm.DB {
	return db.Where("is_checked = ?", false)
}
