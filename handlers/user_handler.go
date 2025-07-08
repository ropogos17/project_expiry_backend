package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"project_expiry_backend/DTO"
	"project_expiry_backend/models"
	"project_expiry_backend/utils"
	"strconv"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	rows, err := h.DB.Query("SELECT * FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Phone, &user.HasAvatar, &user.Password, &user.Salt, &user.EmailVerified, &user.PhoneVerified, &user.TwoFactorEnabled, &user.TwoFactorSecret, &user.LastLogin, &user.LoginAttempts, &user.LockedUntil, &user.PasswordChangedAt, &user.RoleID, &user.IsSuperuser, &user.StatusID, &user.DeletedAt, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	row := h.DB.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err = row.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Phone, &user.HasAvatar, &user.Password, &user.Salt, &user.EmailVerified, &user.PhoneVerified, &user.TwoFactorEnabled, &user.TwoFactorSecret, &user.LastLogin, &user.LoginAttempts, &user.LockedUntil, &user.PasswordChangedAt, &user.RoleID, &user.IsSuperuser, &user.StatusID, &user.DeletedAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var inputUser DTO.UserCreateRequest
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input\n" + err.Error()})
		return
	}

	var user = models.NewUser()
	salt, err := utils.GenerateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate salt"})
		return
	}

	hashedPassword, err := utils.HashPasswordWithSalt(inputUser.Password, salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Name = inputUser.Name
	user.Username = inputUser.Username
	user.Email = inputUser.Email
	user.Password = hashedPassword
	user.Salt = salt

	err = h.DB.QueryRow(
		`INSERT INTO users (name, username, email, phone, password, salt, password_changed_at, role_id, status_id, created_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				RETURNING id`,
		user.Name, user.Username, user.Email, user.Phone, hashedPassword, salt, user.PasswordChangedAt, user.RoleID, user.StatusID, user.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func
