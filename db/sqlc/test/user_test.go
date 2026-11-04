package test

import (
	"context"
	"database/sql"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/utils"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) db.User {
	role := createRandomRole(t)

	arg := db.CreateUserParams{
		Name:        utils.RandomString(10),
		Username:    uuid.New().String(),
		Email:       uuid.New().String() + "@test.com",
		Phone:       sql.NullString{String: "+36301234567", Valid: true},
		HasAvatar:   false,
		Password:    utils.RandomString(32),
		RoleID:      role.ID,
		IsSuperuser: false,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Phone, user.Phone)
	require.Equal(t, arg.RoleID, user.RoleID)
	require.Equal(t, arg.IsSuperuser, user.IsSuperuser)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.False(t, user.DeletedAt.Valid)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByID(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.RoleID, user2.RoleID)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Email, user2.Email)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	arg := db.UpdateUserParams{
		ID:                user1.ID,
		Name:              utils.RandomString(10),
		Username:          uuid.New().String(),
		Email:             uuid.New().String() + "@test.com",
		Phone:             sql.NullString{Valid: false},
		HasAvatar:         false,
		EmailVerified:     false,
		PhoneVerified:     false,
		TwoFactorEnabled:  false,
		TwoFactorSecret:   sql.NullString{Valid: false},
		LastLogin:         sql.NullTime{Valid: false},
		LoginAttempts:     0,
		LockedUntil:       sql.NullTime{Valid: false},
		PasswordChangedAt: sql.NullTime{Valid: false},
		RoleID:            user1.RoleID,
		IsSuperuser:       false,
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, arg.Name, user2.Name)
	require.Equal(t, arg.Username, user2.Username)
	require.Equal(t, arg.Email, user2.Email)
	require.False(t, user2.Phone.Valid)
	require.True(t, user2.UpdatedAt.After(user1.UpdatedAt))
}

func TestSoftDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.SoftDeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)
	require.NoError(t, err)
	require.True(t, user2.DeletedAt.Valid)
}

func TestRestoreUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testQueries.SoftDeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	err = testQueries.RestoreUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)
	require.NoError(t, err)
	require.False(t, user2.DeletedAt.Valid)
}

func TestHardDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.HardDeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserByID(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := db.ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
