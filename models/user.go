package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name      string `json:"name" gorm:"not null" binding:"required"`
	Username  string `json:"username" gorm:"unique;not null" binding:"required"`
	Email     string `json:"email" gorm:"unique;not null" binding:"required,email"`
	Phone     string `json:"phone" gorm:"default:''"`
	HasAvatar bool   `json:"has_avatar" gorm:"default:false"`

	// Authentication fields
	Password          string     `json:"-" gorm:"not null"` // Hidden from JSON
	Salt              string     `json:"-" gorm:"not null"` // Hidden from JSON
	EmailVerified     bool       `json:"email_verified" gorm:"default:false"`
	PhoneVerified     bool       `json:"phone_verified" gorm:"default:false"`
	TwoFactorEnabled  bool       `json:"two_factor_enabled" gorm:"default:false"`
	TwoFactorSecret   *string    `json:"-" gorm:"type:varchar(255)"` // Hidden from JSON
	LastLogin         *time.Time `json:"last_login" gorm:"index"`
	LoginAttempts     int        `json:"login_attempts" gorm:"default:0"`
	LockedUntil       *time.Time `json:"locked_until" gorm:"index"`
	PasswordChangedAt time.Time  `json:"password_changed_at" gorm:"not null"`

	// Relationships
	RoleID      uint `json:"role_id" gorm:"not null;index"`
	Role        Role `json:"role" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IsSuperuser bool `json:"is_superuser" gorm:"default:false;index"`

	StatusID uint       `json:"status_id" gorm:"not null;index"`
	Status   UserStatus `json:"status" gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Reverse relationships
	OwnedStorages         []Storage          `json:"owned_storages,omitempty" gorm:"foreignKey:Owner"`
	ModifiedItems         []Item             `json:"modified_items,omitempty" gorm:"foreignKey:LastModified"`
	ModifiedShoppingItems []ShoppingListItem `json:"modified_shopping_items,omitempty" gorm:"foreignKey:LastModified"`
	SharedFromMe          []Share            `json:"shared_from_me,omitempty" gorm:"foreignKey:SourceUserID"`
	SharedToMe            []Share            `json:"shared_to_me,omitempty" gorm:"foreignKey:TargetUserID"`
	Logs                  []Log              `json:"logs,omitempty" gorm:"foreignKey:UserID"`
}

// TableName overrides the table name used by User to `users`
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM hook - runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.PasswordChangedAt = time.Now()
	return nil
}

// BeforeUpdate GORM hook - runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		u.PasswordChangedAt = time.Now()
	}
	return nil
}
