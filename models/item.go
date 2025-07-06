package models

import "time"

// Item represents an inventory item
type Item struct {
	ID           int        `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	CategoryID   int        `json:"category_id" db:"category_id"`
	Description  *string    `json:"description" db:"description"`
	Quantity     *string    `json:"quantity" db:"quantity"`
	HasImage     bool       `json:"has_image" db:"has_image"`
	IsDeleted    bool       `json:"is_deleted" db:"is_deleted"`
	StoredAt     time.Time  `json:"stored_at" db:"stored_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
	ExpiresAt    time.Time  `json:"expires_at" db:"expires_at"`
	StorageID    int        `json:"storage_id" db:"storage_id"`
	Shelf        *int       `json:"shelf" db:"shelf"`
	LastModified int        `json:"last_modified" db:"last_modified"`
}
