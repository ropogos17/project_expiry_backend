package test

import (
	"context"
	db "project_expiry_backend/db/sqlc"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

func createRandomActionType(t *testing.T) db.ActionType {
	name := uuid.New().String()

	actionType, err := testQueries.CreateActionType(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, actionType)

	require.Equal(t, name, actionType.Name)
	require.NotZero(t, actionType.ID)

	return actionType
}

func TestCreateActionType(t *testing.T) {
	createRandomActionType(t)
}

func TestGetActionTypeByID(t *testing.T) {
	actionType1 := createRandomActionType(t)
	actionType2, err := testQueries.GetActionTypeByID(context.Background(), actionType1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, actionType2)

	require.Equal(t, actionType1.ID, actionType2.ID)
	require.Equal(t, actionType1.Name, actionType2.Name)
}

func TestUpdateActionType(t *testing.T) {
	actionType1 := createRandomActionType(t)
	arg := db.UpdateActionTypeParams{ID: actionType1.ID, Name: uuid.New().String()}

	actionType2, err := testQueries.UpdateActionType(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, actionType2)

	require.Equal(t, actionType1.ID, actionType2.ID)
	require.Equal(t, arg.Name, actionType2.Name)
}

func TestHardDeleteActionType(t *testing.T) {
	actionType1 := createRandomActionType(t)
	rows, err := testQueries.HardDeleteActionType(context.Background(), actionType1.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)

	actionType2, err := testQueries.GetActionTypeByID(context.Background(), actionType1.ID)
	require.Error(t, err)
	require.Empty(t, actionType2)

	rows, err = testQueries.HardDeleteActionType(context.Background(), actionType1.ID)
	require.NoError(t, err)
	require.Equal(t, int64(0), rows)
}

func TestListActionTypes(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomActionType(t)
	}

	arg := db.ListActionTypesParams{Limit: 5, Offset: 5}

	actionTypes, err := testQueries.ListActionTypes(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, actionTypes, 5)

	for _, at := range actionTypes {
		require.NotEmpty(t, at)
	}
}

func TestActionTypeExistsByName(t *testing.T) {
	actionType := createRandomActionType(t)

	exists, err := testQueries.ActionTypeExistsByName(context.Background(), actionType.Name)
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = testQueries.ActionTypeExistsByName(context.Background(), uuid.New().String())
	require.NoError(t, err)
	require.False(t, exists)
}
