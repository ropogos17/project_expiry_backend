package DTO

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"` //binding:"required,matches=^(\\+36)?\\s?(20|30|31|70|1|2[1-9]|3[1-9]|4[1-9]|5[1-9]|6[1-9])\\s?\\d{3}\\s?\\d{3,4}$"
	Password string `json:"password" binding:"required,min=8"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	RoleID   *uint   `json:"role_id,omitempty"`
	StatusID *uint   `json:"status_id,omitempty"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
