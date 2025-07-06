package models

// ActionType represents the type of action performed
type ActionType struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
