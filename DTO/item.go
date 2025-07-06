package DTO

// ItemCreateRequest represents the request body for creating an item
type ItemCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	CategoryID  int     `json:"category_id" binding:"required"`
	Description *string `json:"description"`
	Quantity    *string `json:"quantity"`
	ExpiresAt   string  `json:"expires_at" binding:"required"` // ISO format
	StorageID   int     `json:"storage_id" binding:"required"`
	Shelf       *int    `json:"shelf"`
}

// ItemUpdateRequest represents the request body for updating an item
type ItemUpdateRequest struct {
	Name        *string `json:"name"`
	CategoryID  *int    `json:"category_id"`
	Description *string `json:"description"`
	Quantity    *string `json:"quantity"`
	ExpiresAt   *string `json:"expires_at"` // ISO format
	StorageID   *int    `json:"storage_id"`
	Shelf       *int    `json:"shelf"`
}
