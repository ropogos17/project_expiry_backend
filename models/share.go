package models

import (
	"gorm.io/gorm"
	"time"
)

// Share represents the share table
type Share struct {
	gorm.Model
	SourceUserID uint    `json:"source_user_id" gorm:"not null;index;uniqueIndex:idx_share_unique"`
	SourceUser   User    `json:"source_user" gorm:"foreignKey:SourceUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StorageID    uint    `json:"storage_id" gorm:"not null;index;uniqueIndex:idx_share_unique"`
	Storage      Storage `json:"storage" gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TargetUserID uint    `json:"target_user_id" gorm:"not null;index;uniqueIndex:idx_share_unique"`
	TargetUser   User    `json:"target_user" gorm:"foreignKey:TargetUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CanShare     bool    `json:"can_share" gorm:"default:false"`

	// Timestamps
	SharedAt time.Time `json:"shared_at" gorm:"not null;default:CURRENT_TIMESTAMP;index"`
}

// BeforeCreate GORM hook - runs before creating a share
func (s *Share) BeforeCreate(tx *gorm.DB) error {
	s.SharedAt = time.Now()
	return nil
}
