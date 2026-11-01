package test

import (
	"context"
	"github.com/google/uuid"
	db "project_expiry_backend/db/sqlc"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T) db.Category {
	name := uuid.New().String()

	category, err := testQueries.CreateCategory(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.Equal(t, name, category.Name)
	require.NotZero(t, category.ID)

	return category
}

func TestCreateCategory(t *testing.T) {
	createRandomCategory(t)
}

func TestGetCategoryByID(t *testing.T) {
	category1 := createRandomCategory(t)
	category2, err := testQueries.GetCategoryByID(context.Background(), category1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.Equal(t, category1.ID, category2.ID)
	require.Equal(t, category1.Name, category2.Name)
}

func TestUpdateCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	arg := db.UpdateCategoryParams{ID: category1.ID, Name: uuid.New().String()}

	category2, err := testQueries.UpdateCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.Equal(t, category1.ID, category2.ID)
	require.Equal(t, arg.Name, category2.Name)
}

func TestHardDeleteCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	err := testQueries.HardDeleteCategory(context.Background(), category1.ID)
	require.NoError(t, err)

	category2, err := testQueries.GetCategoryByID(context.Background(), category1.ID)
	require.Error(t, err)
	require.Empty(t, category2)
}

func TestListCategories(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCategory(t)
	}

	arg := db.ListCategoriesParams{Limit: 5, Offset: 5}

	categories, err := testQueries.ListCategories(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, categories, 5)

	for _, c := range categories {
		require.NotEmpty(t, c)
	}
}

func TestCategoryExistsByName(t *testing.T) {
	category := createRandomCategory(t)

	exists, err := testQueries.CategoryExistsByName(context.Background(), category.Name)
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = testQueries.CategoryExistsByName(context.Background(), uuid.New().String())
	require.NoError(t, err)
	require.False(t, exists)
}
