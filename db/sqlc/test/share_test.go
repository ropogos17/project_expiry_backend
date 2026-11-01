package test

import (
	"context"
	db "project_expiry_backend/db/sqlc"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomShare(t *testing.T) db.Share {
	sourceUser := createRandomUser(t)
	targetUser := createRandomUser(t)
	storage := createRandomStorage(t)

	arg := db.CreateShareParams{
		SourceUserID: sourceUser.ID,
		StorageID:    storage.ID,
		TargetUserID: targetUser.ID,
		CanShare:     false,
	}

	share, err := testQueries.CreateShare(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, share)

	require.Equal(t, arg.SourceUserID, share.SourceUserID)
	require.Equal(t, arg.StorageID, share.StorageID)
	require.Equal(t, arg.TargetUserID, share.TargetUserID)
	require.Equal(t, arg.CanShare, share.CanShare)
	require.NotZero(t, share.ID)

	return share
}

func TestCreateShare(t *testing.T) {
	createRandomShare(t)
}

func TestGetShareByID(t *testing.T) {
	share1 := createRandomShare(t)
	share2, err := testQueries.GetShareByID(context.Background(), share1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, share2)

	require.Equal(t, share1.ID, share2.ID)
	require.Equal(t, share1.SourceUserID, share2.SourceUserID)
	require.Equal(t, share1.TargetUserID, share2.TargetUserID)
	require.Equal(t, share1.StorageID, share2.StorageID)
}

func TestGetShareByStorageAndTargetUser(t *testing.T) {
	share1 := createRandomShare(t)

	arg := db.GetShareByStorageAndTargetUserParams{
		StorageID:    share1.StorageID,
		TargetUserID: share1.TargetUserID,
	}

	share2, err := testQueries.GetShareByStorageAndTargetUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, share2)

	require.Equal(t, share1.ID, share2.ID)
}

func TestUpdateShare(t *testing.T) {
	share1 := createRandomShare(t)

	arg := db.UpdateShareParams{
		ID:           share1.ID,
		SourceUserID: share1.SourceUserID,
		StorageID:    share1.StorageID,
		TargetUserID: share1.TargetUserID,
		CanShare:     true,
		SharedAt:     time.Now(),
	}

	share2, err := testQueries.UpdateShare(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, share2)

	require.Equal(t, share1.ID, share2.ID)
	require.True(t, share2.CanShare)
}

func TestHardDeleteShare(t *testing.T) {
	share1 := createRandomShare(t)
	err := testQueries.HardDeleteShare(context.Background(), share1.ID)
	require.NoError(t, err)

	share2, err := testQueries.GetShareByID(context.Background(), share1.ID)
	require.Error(t, err)
	require.Empty(t, share2)
}

func TestListShares(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomShare(t)
	}

	arg := db.ListSharesParams{Limit: 5, Offset: 5}

	shares, err := testQueries.ListShares(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, shares, 5)

	for _, s := range shares {
		require.NotEmpty(t, s)
	}
}

func TestListSharesByUser(t *testing.T) {
	sourceUser := createRandomUser(t)
	targetUser := createRandomUser(t)

	for i := 0; i < 3; i++ {
		storage := createRandomStorage(t)
		arg := db.CreateShareParams{
			SourceUserID: sourceUser.ID,
			StorageID:    storage.ID,
			TargetUserID: targetUser.ID,
			CanShare:     false,
		}
		_, err := testQueries.CreateShare(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := db.ListSharesByUserParams{SourceUserID: sourceUser.ID, Limit: 10, Offset: 0}

	shares, err := testQueries.ListSharesByUser(context.Background(), arg)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(shares), 3)

	for _, s := range shares {
		require.Equal(t, sourceUser.ID, s.SourceUserID)
	}
}
