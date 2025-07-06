package models

// UserStatus represents the status table
type Status struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
