package models

import "gorm.io/gorm"

// Log represents the log table
type Log struct {
	gorm.Model
	ActionTypeID *uint       `json:"action_type_id" gorm:"index"`
	ActionType   *ActionType `json:"action_type,omitempty" gorm:"foreignKey:ActionTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Message      *string     `json:"message" gorm:"type:text"`
	Data         *string     `json:"data" gorm:"type:text"`
	UserID       uint        `json:"user_id" gorm:"not null;index"`
	User         User        `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
