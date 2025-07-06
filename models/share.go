package models

import "time"

// Share represents a storage sharing relationship
type Share struct {
	ID           int       `json:"id" db:"id"`
	SourceUserID int       `json:"source_user_id" db:"source_user_id"`
	StorageID    int       `json:"storage_id" db:"storage_id"`
	TargetUserID int       `json:"target_user_id" db:"target_user_id"`
	CanShare     bool      `json:"can_share" db:"can_share"`
	SharedAt     time.Time `json:"shared_at" db:"shared_at"`
}
