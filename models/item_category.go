package models

// ItemCategory represents an item's category
type ItemCategory struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
