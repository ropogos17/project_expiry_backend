package models

// Role represents the role table
type Role struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
