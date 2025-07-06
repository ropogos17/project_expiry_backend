package models

// Storage represents a storage location
type Storage struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	Location    *string `json:"location" db:"location"`
	Shelves     *int    `json:"shelves" db:"shelves"`
	HasImage    bool    `json:"has_image" db:"has_image"`
	Owner       int     `json:"owner" db:"owner"`
}
