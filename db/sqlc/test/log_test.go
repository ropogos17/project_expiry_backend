package test

import (
	"context"
	"database/sql"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/utils"
	"testing"

	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/require"
)

func createRandomLog(t *testing.T) db.Log {
	actionType := createRandomActionType(t)
	user := createRandomUser(t)

	arg := db.CreateLogParams{
		ActionTypeID: actionType.ID,
		Message:      sql.NullString{String: utils.RandomString(50), Valid: true},
		Data:         pqtype.NullRawMessage{Valid: false},
		UserID:       user.ID,
	}

	log, err := testQueries.CreateLog(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, log)

	require.Equal(t, arg.ActionTypeID, log.ActionTypeID)
	require.Equal(t, arg.Message, log.Message)
	require.Equal(t, arg.UserID, log.UserID)
	require.NotZero(t, log.ID)
	require.NotZero(t, log.CreatedAt)

	return log
}

func TestCreateLog(t *testing.T) {
	createRandomLog(t)
}

func TestGetLogByID(t *testing.T) {
	log1 := createRandomLog(t)
	log2, err := testQueries.GetLogByID(context.Background(), log1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, log2)

	require.Equal(t, log1.ID, log2.ID)
	require.Equal(t, log1.ActionTypeID, log2.ActionTypeID)
	require.Equal(t, log1.Message, log2.Message)
	require.Equal(t, log1.UserID, log2.UserID)
}

func TestUpdateLog(t *testing.T) {
	log1 := createRandomLog(t)
	newActionType := createRandomActionType(t)

	arg := db.UpdateLogParams{
		ID:           log1.ID,
		ActionTypeID: newActionType.ID,
		Message:      sql.NullString{String: utils.RandomString(50), Valid: true},
		Data:         pqtype.NullRawMessage{Valid: false},
	}

	log2, err := testQueries.UpdateLog(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, log2)

	require.Equal(t, log1.ID, log2.ID)
	require.Equal(t, arg.ActionTypeID, log2.ActionTypeID)
	require.Equal(t, arg.Message, log2.Message)
}

func TestHardDeleteLog(t *testing.T) {
	log1 := createRandomLog(t)
	err := testQueries.HardDeleteLog(context.Background(), log1.ID)
	require.NoError(t, err)

	log2, err := testQueries.GetLogByID(context.Background(), log1.ID)
	require.Error(t, err)
	require.Empty(t, log2)
}

func TestListLogs(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomLog(t)
	}

	arg := db.ListLogsParams{Limit: 5, Offset: 5}

	logs, err := testQueries.ListLogs(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, logs, 5)

	for _, l := range logs {
		require.NotEmpty(t, l)
	}
}
