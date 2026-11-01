package test

import (
	"context"
	db "project_expiry_backend/db/sqlc"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomRefreshToken(t *testing.T) db.RefreshToken {
	user := createRandomUser(t)

	arg := db.CreateRefreshTokenParams{
		UserID:    user.ID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	token, err := testQueries.CreateRefreshToken(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.Equal(t, arg.UserID, token.UserID)
	require.NotZero(t, token.ID)
	require.False(t, token.DeletedAt.Valid)

	return token
}

func TestCreateRefreshToken(t *testing.T) {
	createRandomRefreshToken(t)
}

func TestGetRefreshTokenByID(t *testing.T) {
	token1 := createRandomRefreshToken(t)
	token2, err := testQueries.GetRefreshTokenByID(context.Background(), token1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, token2)

	require.Equal(t, token1.ID, token2.ID)
	require.Equal(t, token1.UserID, token2.UserID)
}

func TestSoftDeleteRefreshToken(t *testing.T) {
	token1 := createRandomRefreshToken(t)
	err := testQueries.SoftDeleteRefreshToken(context.Background(), token1.ID)
	require.NoError(t, err)

	// GetRefreshTokenByID csak deleted_at IS NULL rekordot ad vissza
	token2, err := testQueries.GetRefreshTokenByID(context.Background(), token1.ID)
	require.Error(t, err)
	require.Empty(t, token2)
}

func TestRestoreRefreshToken(t *testing.T) {
	token1 := createRandomRefreshToken(t)

	err := testQueries.SoftDeleteRefreshToken(context.Background(), token1.ID)
	require.NoError(t, err)

	err = testQueries.RestoreRefreshToken(context.Background(), token1.ID)
	require.NoError(t, err)

	token2, err := testQueries.GetRefreshTokenByID(context.Background(), token1.ID)
	require.NoError(t, err)
	require.False(t, token2.DeletedAt.Valid)
}

func TestHardDeleteRefreshToken(t *testing.T) {
	token1 := createRandomRefreshToken(t)
	err := testQueries.HardDeleteRefreshToken(context.Background(), token1.ID)
	require.NoError(t, err)

	token2, err := testQueries.GetRefreshTokenByID(context.Background(), token1.ID)
	require.Error(t, err)
	require.Empty(t, token2)
}

func TestListRefreshTokens(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomRefreshToken(t)
	}

	arg := db.ListRefreshTokensParams{Limit: 5, Offset: 5}

	tokens, err := testQueries.ListRefreshTokens(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tokens, 5)

	for _, tok := range tokens {
		require.NotEmpty(t, tok)
	}
}

func TestListRefreshTokensByUser(t *testing.T) {
	user := createRandomUser(t)

	for i := 0; i < 3; i++ {
		arg := db.CreateRefreshTokenParams{
			UserID:    user.ID,
			IssuedAt:  time.Now(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		_, err := testQueries.CreateRefreshToken(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := db.ListRefreshTokensByUserParams{UserID: user.ID, Limit: 10, Offset: 0}

	tokens, err := testQueries.ListRefreshTokensByUser(context.Background(), arg)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(tokens), 3)

	for _, tok := range tokens {
		require.Equal(t, user.ID, tok.UserID)
	}
}

func TestUpdateRefreshToken(t *testing.T) {
	token1 := createRandomRefreshToken(t)
	user := createRandomUser(t)

	arg := db.UpdateRefreshTokenParams{
		ID:        token1.ID,
		UserID:    user.ID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	token2, err := testQueries.UpdateRefreshToken(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, token2)

	require.Equal(t, token1.ID, token2.ID)
	require.Equal(t, arg.UserID, token2.UserID)
	require.True(t, token2.ExpiresAt.After(token1.ExpiresAt))
}
