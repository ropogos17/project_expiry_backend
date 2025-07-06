package models

// Announcement represents a system announcement
type Announcement struct {
	ID           int     `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	Description  *string `json:"description" db:"description"`
	ExternalLink *string `json:"external_link" db:"external_link"`
	HasRead      bool    `json:"has_read" db:"has_read"`
}
