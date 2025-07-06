package models

import (
	"gorm.io/gorm"
	"time"
)

// Item represents the item table
type Item struct {
	gorm.Model
	Name        string       `json:"name" gorm:"not null;index" binding:"required"`
	CategoryID  uint         `json:"category_id" gorm:"not null;index"`
	Category    ItemCategory `json:"category" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Description *string      `json:"description" gorm:"type:text"`
	Quantity    *string      `json:"quantity"`
	HasImage    bool         `json:"has_image" gorm:"default:false"`
	IsDeleted   bool         `json:"is_deleted" gorm:"default:false;index"`

	// Timestamps
	StoredAt  time.Time `json:"stored_at" gorm:"not null;default:CURRENT_TIMESTAMP;index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`

	// Storage relationship
	StorageID uint    `json:"storage_id" gorm:"not null;index"`
	Storage   Storage `json:"storage" gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Shelf     *int    `json:"shelf"`

	// User relationship
	LastModified uint `json:"last_modified" gorm:"not null;index"`
	ModifiedBy   User `json:"modified_by" gorm:"foreignKey:LastModified;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// BeforeCreate GORM hook - runs before creating an item
func (i *Item) BeforeCreate(tx *gorm.DB) error {
	i.StoredAt = time.Now()
	return nil
}

// Scopes for common queries
func (Item) NotDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("is_deleted = ?", false)
}

func (Item) ByStorage(storageID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("storage_id = ?", storageID)
	}
}

func (Item) ExpiringBefore(date time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("expires_at < ?", date)
	}
}
