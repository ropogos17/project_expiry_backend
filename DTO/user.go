package DTO

import (
	"project_expiry_backend/models"
)

// UserCreateRequest represents the request body for creating a user
type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

// UserUpdateRequest represents the request body for updating a user
type UserUpdateRequest struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"` //binding:"required,matches=^(\\+36)?\\s?(20|30|31|70|1|2[1-9]|3[1-9]|4[1-9]|5[1-9]|6[1-9])\\s?\\d{3}\\s?\\d{3,4}$"
	RoleID   *int    `json:"role_id"`
}

// UserLoginRequest represents the request body for user login
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponse represents the response body for user login
type UserLoginResponse struct {
	User         models.User `json:"user"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
}
