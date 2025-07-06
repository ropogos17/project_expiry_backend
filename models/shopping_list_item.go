package models

import "time"

// ShoppingListItem represents an item in a shopping list
type ShoppingListItem struct {
	ID           int        `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Description  *string    `json:"description" db:"description"`
	Quantity     *string    `json:"quantity" db:"quantity"`
	IsChecked    bool       `json:"is_checked" db:"is_checked"`
	IsDeleted    bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
	LastModified int        `json:"last_modified" db:"last_modified"`
}
