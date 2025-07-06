package DTO

import "time"

// CreateItemRequest represents the request body for creating an item
type CreateItemRequest struct {
	Name        string    `json:"name" binding:"required"`
	CategoryID  uint      `json:"category_id" binding:"required"`
	Description *string   `json:"description"`
	Quantity    *string   `json:"quantity"`
	StorageID   uint      `json:"storage_id" binding:"required"`
	Shelf       *int      `json:"shelf"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
}

// UpdateItemRequest represents the request body for updating an item
type UpdateItemRequest struct {
	Name        *string    `json:"name,omitempty"`
	CategoryID  *uint      `json:"category_id,omitempty"`
	Description *string    `json:"description,omitempty"`
	Quantity    *string    `json:"quantity,omitempty"`
	StorageID   *uint      `json:"storage_id,omitempty"`
	Shelf       *int       `json:"shelf,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}
