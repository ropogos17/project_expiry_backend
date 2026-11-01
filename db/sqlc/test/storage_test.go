package test

import (
	"context"
	"database/sql"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomStorage(t *testing.T) db.Storage {
	user := createRandomUser(t)

	arg := db.CreateStorageParams{
		Name:        utils.RandomString(10),
		Description: sql.NullString{String: utils.RandomString(50), Valid: true},
		Location:    sql.NullString{Valid: false},
		ShelfCount:  int32(utils.RandomInt(1, 20)),
		HasImage:    false,
		Owner:       user.ID,
	}

	storage, err := testQueries.CreateStorage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, storage)

	require.Equal(t, arg.Name, storage.Name)
	require.Equal(t, arg.Description, storage.Description)
	require.False(t, storage.Location.Valid)
	require.Equal(t, arg.ShelfCount, storage.ShelfCount)
	require.Equal(t, arg.Owner, storage.Owner)
	require.NotZero(t, storage.ID)
	require.False(t, storage.DeletedAt.Valid)

	return storage
}

func TestCreateStorage(t *testing.T) {
	createRandomStorage(t)
}

func TestGetStorageByID(t *testing.T) {
	storage1 := createRandomStorage(t)
	storage2, err := testQueries.GetStorageByID(context.Background(), storage1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, storage2)

	require.Equal(t, storage1.ID, storage2.ID)
	require.Equal(t, storage1.Name, storage2.Name)
	require.Equal(t, storage1.Owner, storage2.Owner)
}

func TestUpdateStorage(t *testing.T) {
	storage1 := createRandomStorage(t)

	arg := db.UpdateStorageParams{
		ID:          storage1.ID,
		Name:        utils.RandomString(10),
		Description: sql.NullString{String: utils.RandomString(50), Valid: true},
		Location:    sql.NullString{String: utils.RandomString(20), Valid: true},
		ShelfCount:  int32(utils.RandomInt(1, 20)),
		HasImage:    false,
		Owner:       storage1.Owner,
	}

	storage2, err := testQueries.UpdateStorage(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, storage2)

	require.Equal(t, storage1.ID, storage2.ID)
	require.Equal(t, arg.Name, storage2.Name)
	require.Equal(t, arg.Location, storage2.Location)
	require.True(t, storage2.UpdatedAt.After(storage1.UpdatedAt))
}

func TestDeleteStorage(t *testing.T) {
	storage1 := createRandomStorage(t)
	err := testQueries.DeleteStorage(context.Background(), storage1.ID)
	require.NoError(t, err)

	storage2, err := testQueries.GetStorageByID(context.Background(), storage1.ID)
	require.NoError(t, err)
	require.True(t, storage2.DeletedAt.Valid)
}

func TestRestoreStorage(t *testing.T) {
	storage1 := createRandomStorage(t)

	err := testQueries.DeleteStorage(context.Background(), storage1.ID)
	require.NoError(t, err)

	err = testQueries.RestoreStorage(context.Background(), storage1.ID)
	require.NoError(t, err)

	storage2, err := testQueries.GetStorageByID(context.Background(), storage1.ID)
	require.NoError(t, err)
	require.False(t, storage2.DeletedAt.Valid)
}

func TestHardDeleteStorage(t *testing.T) {
	storage1 := createRandomStorage(t)
	err := testQueries.HardDeleteStorage(context.Background(), storage1.ID)
	require.NoError(t, err)

	storage2, err := testQueries.GetStorageByID(context.Background(), storage1.ID)
	require.Error(t, err)
	require.Empty(t, storage2)
}

func TestListStorages(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomStorage(t)
	}

	arg := db.ListStoragesParams{Limit: 5, Offset: 5}

	storages, err := testQueries.ListStorages(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, storages, 5)

	for _, s := range storages {
		require.NotEmpty(t, s)
	}
}

func TestListStoragesByOwner(t *testing.T) {
	user := createRandomUser(t)

	for i := 0; i < 3; i++ {
		arg := db.CreateStorageParams{
			Name:       utils.RandomString(10),
			ShelfCount: int32(utils.RandomInt(1, 10)),
			HasImage:   false,
			Owner:      user.ID,
		}
		_, err := testQueries.CreateStorage(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := db.ListStoragesByOwnerParams{Limit: 10, Offset: 0, Owner: user.ID}

	storages, err := testQueries.ListStoragesByOwner(context.Background(), arg)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(storages), 3)

	for _, s := range storages {
		require.Equal(t, user.ID, s.Owner)
		require.False(t, s.DeletedAt.Valid)
	}
}
