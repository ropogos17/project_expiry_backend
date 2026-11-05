package test

import (
	"context"
	"database/sql"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomItem(t *testing.T) db.Item {
	category := createRandomCategory(t)
	storage := createRandomStorage(t)
	user := createRandomUser(t)

	arg := db.CreateItemParams{
		Name:           utils.RandomString(10),
		CategoryID:     category.ID,
		Description:    sql.NullString{String: utils.RandomString(50), Valid: true},
		Quantity:       sql.NullString{Valid: false},
		HasImage:       false,
		StoredAt:       time.Now(),
		ExpiresAt:      time.Now().Add(30 * 24 * time.Hour),
		StorageID:      storage.ID,
		Shelf:          int32(utils.RandomInt(1, 10)),
		LastModifiedBy: user.ID,
	}

	item, err := testQueries.CreateItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, item)

	require.Equal(t, arg.Name, item.Name)
	require.Equal(t, arg.CategoryID, item.CategoryID)
	require.Equal(t, arg.StorageID, item.StorageID)
	require.Equal(t, arg.LastModifiedBy, item.LastModifiedBy)
	require.NotZero(t, item.ID)
	require.False(t, item.DeletedAt.Valid)

	return item
}

func TestCreateItem(t *testing.T) {
	createRandomItem(t)
}

func TestGetItemByID(t *testing.T) {
	item1 := createRandomItem(t)
	item2, err := testQueries.GetItemByID(context.Background(), item1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, item1.ID, item2.ID)
	require.Equal(t, item1.Name, item2.Name)
	require.Equal(t, item1.StorageID, item2.StorageID)
}

func TestUpdateItem(t *testing.T) {
	item1 := createRandomItem(t)

	arg := db.UpdateItemParams{
		ID:             item1.ID,
		Name:           utils.RandomString(10),
		CategoryID:     item1.CategoryID,
		Description:    sql.NullString{String: utils.RandomString(50), Valid: true},
		Quantity:       sql.NullString{String: "5 kg", Valid: true},
		HasImage:       false,
		StoredAt:       item1.StoredAt,
		ExpiresAt:      time.Now().Add(60 * 24 * time.Hour),
		StorageID:      item1.StorageID,
		Shelf:          item1.Shelf,
		LastModifiedBy: item1.LastModifiedBy,
	}

	item2, err := testQueries.UpdateItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, item1.ID, item2.ID)
	require.Equal(t, arg.Name, item2.Name)
	require.Equal(t, arg.Quantity, item2.Quantity)
	require.True(t, item2.UpdatedAt.After(item1.UpdatedAt))
}

func TestDeleteItem(t *testing.T) {
	item1 := createRandomItem(t)
	err := testQueries.DeleteItem(context.Background(), item1.ID)
	require.NoError(t, err)

	item2, err := testQueries.GetItemByID(context.Background(), item1.ID)
	require.NoError(t, err)
	require.True(t, item2.DeletedAt.Valid)
}

func TestRestoreItem(t *testing.T) {
	item1 := createRandomItem(t)

	err := testQueries.DeleteItem(context.Background(), item1.ID)
	require.NoError(t, err)

	err = testQueries.RestoreItem(context.Background(), item1.ID)
	require.NoError(t, err)

	item2, err := testQueries.GetItemByID(context.Background(), item1.ID)
	require.NoError(t, err)
	require.False(t, item2.DeletedAt.Valid)
}

func TestHardDeleteItem(t *testing.T) {
	item1 := createRandomItem(t)
	rows, err := testQueries.HardDeleteItem(context.Background(), item1.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)

	item2, err := testQueries.GetItemByID(context.Background(), item1.ID)
	require.Error(t, err)
	require.Empty(t, item2)

	rows, err = testQueries.HardDeleteItem(context.Background(), item1.ID)
	require.NoError(t, err)
	require.Equal(t, int64(0), rows)
}

func TestListItems(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomItem(t)
	}

	arg := db.ListItemsParams{Limit: 5, Offset: 5}

	items, err := testQueries.ListItems(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, items, 5)
}

func TestListItemsByStorage(t *testing.T) {
	item1 := createRandomItem(t)

	arg := db.ListItemsByStorageParams{StorageID: item1.StorageID, Limit: 10, Offset: 0}

	items, err := testQueries.ListItemsByStorage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)

	for _, item := range items {
		require.Equal(t, item1.StorageID, item.StorageID)
		require.False(t, item.DeletedAt.Valid)
	}
}

func TestListItemsByUser(t *testing.T) {
	storage := createRandomStorage(t)

	arg := db.CreateItemParams{
		Name:           utils.RandomString(10),
		CategoryID:     createRandomCategory(t).ID,
		HasImage:       false,
		StoredAt:       time.Now(),
		ExpiresAt:      time.Now().Add(30 * 24 * time.Hour),
		StorageID:      storage.ID,
		Shelf:          1,
		LastModifiedBy: storage.Owner,
	}
	_, err := testQueries.CreateItem(context.Background(), arg)
	require.NoError(t, err)

	listArg := db.ListItemsByUserParams{Owner: storage.Owner, Limit: 10, Offset: 0}

	items, err := testQueries.ListItemsByUser(context.Background(), listArg)
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestListActiveItemsByStorage(t *testing.T) {
	item1 := createRandomItem(t) // expires_at = now + 30 nap

	arg := db.ListActiveItemsByStorageParams{
		StorageID: item1.StorageID,
		ExpiresAt: time.Now(),
		Limit:     10,
		Offset:    0,
	}

	items, err := testQueries.ListActiveItemsByStorage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)

	for _, item := range items {
		require.True(t, item.ExpiresAt.After(time.Now()))
		require.False(t, item.DeletedAt.Valid)
	}
}

func TestListExpiredItemsByStorage(t *testing.T) {
	category := createRandomCategory(t)
	storage := createRandomStorage(t)
	user := createRandomUser(t)

	arg := db.CreateItemParams{
		Name:           utils.RandomString(10),
		CategoryID:     category.ID,
		HasImage:       false,
		StoredAt:       time.Now().Add(-60 * 24 * time.Hour),
		ExpiresAt:      time.Now().Add(-1 * time.Hour), // már lejárt
		StorageID:      storage.ID,
		Shelf:          1,
		LastModifiedBy: user.ID,
	}
	_, err := testQueries.CreateItem(context.Background(), arg)
	require.NoError(t, err)

	listArg := db.ListExpiredItemsByStorageParams{
		StorageID: storage.ID,
		ExpiresAt: time.Now(),
		Limit:     10,
		Offset:    0,
	}

	items, err := testQueries.ListExpiredItemsByStorage(context.Background(), listArg)
	require.NoError(t, err)
	require.NotEmpty(t, items)

	for _, item := range items {
		require.True(t, item.ExpiresAt.Before(time.Now()))
	}
}

func TestListActiveItemsByUser(t *testing.T) {
	storage := createRandomStorage(t)

	_, err := testQueries.CreateItem(context.Background(), db.CreateItemParams{
		Name:           utils.RandomString(10),
		CategoryID:     createRandomCategory(t).ID,
		HasImage:       false,
		StoredAt:       time.Now(),
		ExpiresAt:      time.Now().Add(30 * 24 * time.Hour),
		StorageID:      storage.ID,
		Shelf:          1,
		LastModifiedBy: storage.Owner,
	})
	require.NoError(t, err)

	arg := db.ListActiveItemsByUserParams{
		Owner:     storage.Owner,
		ExpiresAt: time.Now(),
		Limit:     10,
		Offset:    0,
	}

	items, err := testQueries.ListActiveItemsByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)

	for _, item := range items {
		require.True(t, item.ExpiresAt.After(time.Now()))
		require.False(t, item.DeletedAt.Valid)
	}
}

func TestListDeletedItemsByStorage(t *testing.T) {
	item1 := createRandomItem(t)

	err := testQueries.DeleteItem(context.Background(), item1.ID)
	require.NoError(t, err)

	arg := db.ListDeletedItemsByStorageParams{
		StorageID: item1.StorageID,
		Limit:     10,
		Offset:    0,
	}

	items, err := testQueries.ListDeletedItemsByStorage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)

	for _, item := range items {
		require.True(t, item.DeletedAt.Valid)
	}
}

func TestListDeletedItemsByUser(t *testing.T) {
	storage := createRandomStorage(t)

	item, err := testQueries.CreateItem(context.Background(), db.CreateItemParams{
		Name:           utils.RandomString(10),
		CategoryID:     createRandomCategory(t).ID,
		HasImage:       false,
		StoredAt:       time.Now(),
		ExpiresAt:      time.Now().Add(30 * 24 * time.Hour),
		StorageID:      storage.ID,
		Shelf:          1,
		LastModifiedBy: storage.Owner,
	})
	require.NoError(t, err)

	err = testQueries.DeleteItem(context.Background(), item.ID)
	require.NoError(t, err)

	arg := db.ListDeletedItemsByUserParams{
		Owner:  storage.Owner,
		Limit:  10,
		Offset: 0,
	}

	items, err := testQueries.ListDeletedItemsByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)

	for _, item := range items {
		require.True(t, item.DeletedAt.Valid)
	}
}

func TestListExpiredItemsByUser(t *testing.T) {
	storage := createRandomStorage(t)

	_, err := testQueries.CreateItem(context.Background(), db.CreateItemParams{
		Name:           utils.RandomString(10),
		CategoryID:     createRandomCategory(t).ID,
		HasImage:       false,
		StoredAt:       time.Now().Add(-60 * 24 * time.Hour),
		ExpiresAt:      time.Now().Add(-1 * time.Hour),
		StorageID:      storage.ID,
		Shelf:          1,
		LastModifiedBy: storage.Owner,
	})
	require.NoError(t, err)

	arg := db.ListExpiredItemsByUserParams{
		Owner:     storage.Owner,
		ExpiresAt: time.Now(),
		Limit:     10,
		Offset:    0,
	}

	items, err := testQueries.ListExpiredItemsByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)

	for _, item := range items {
		require.True(t, item.ExpiresAt.Before(time.Now()))
		require.False(t, item.DeletedAt.Valid)
	}
}
