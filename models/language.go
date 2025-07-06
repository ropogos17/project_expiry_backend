package models

import "gorm.io/gorm"

// Language represents the language table
type Language struct {
	gorm.Model
	EN string `json:"en" gorm:"not null;type:text" binding:"required"` // English
	HU string `json:"hu" gorm:"not null;type:text" binding:"required"` // Hungarian
	RO string `json:"ro" gorm:"not null;type:text" binding:"required"` // Romanian
	ES string `json:"es" gorm:"not null;type:text" binding:"required"` // Spanish
	ZH string `json:"zh" gorm:"not null;type:text" binding:"required"` // Chinese (Mandarin)
	FR string `json:"fr" gorm:"not null;type:text" binding:"required"` // French
	DE string `json:"de" gorm:"not null;type:text" binding:"required"` // German
	PT string `json:"pt" gorm:"not null;type:text" binding:"required"` // Portuguese
	RU string `json:"ru" gorm:"not null;type:text" binding:"required"` // Russian
	UA string `json:"ua" gorm:"not null;type:text" binding:"required"` // Ukrainian
	AR string `json:"ar" gorm:"not null;type:text" binding:"required"` // Arabic
}
