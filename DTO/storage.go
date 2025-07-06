package DTO

// StorageCreateRequest represents the request body for creating a storage
type StorageCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	Shelves     *int    `json:"shelves"`
}

// StorageUpdateRequest represents the request body for updating a storage
type StorageUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	Shelves     *int    `json:"shelves"`
}
