package models

import "gorm.io/gorm"

// Storage represents the storage table
type Storage struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null;index" binding:"required"`
	Description *string `json:"description" gorm:"type:text"`
	Location    *string `json:"location" gorm:"index"`
	Shelves     *int    `json:"shelves"`
	HasImage    bool    `json:"has_image" gorm:"default:false"`
	Owner       uint    `json:"owner" gorm:"not null;index"`
	OwnerUser   User    `json:"owner_user" gorm:"foreignKey:Owner;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Reverse relationships
	Items  []Item  `json:"items,omitempty" gorm:"foreignKey:StorageID"`
	Shares []Share `json:"shares,omitempty" gorm:"foreignKey:StorageID"`
}

// Scopes for common queries
func (Storage) ByOwner(ownerID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("owner = ?", ownerID)
	}
}
