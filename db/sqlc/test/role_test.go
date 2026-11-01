package test

import (
	"context"
	"database/sql"
	db "project_expiry_backend/db/sqlc"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomRole(t *testing.T) db.Role {
	arg := uuid.New().String()

	role, err := testQueries.CreateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role)

	require.Equal(t, arg, role.Name)
	require.NotZero(t, role.ID)

	return role
}

func TestCreateRole(t *testing.T) {
	createRandomRole(t)
}

func TestGetRoleByID(t *testing.T) {
	role1 := createRandomRole(t)
	role2, err := testQueries.GetRoleByID(context.Background(), role1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.ID, role2.ID)
	require.Equal(t, role1.Name, role2.Name)
}

func TestUpdateRole(t *testing.T) {
	role1 := createRandomRole(t)
	arg := db.UpdateRoleParams{ID: role1.ID, Name: uuid.New().String()}

	role2, err := testQueries.UpdateRole(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.ID, role2.ID)
	require.Equal(t, arg.Name, role2.Name)
}

func TestHardDeleteRole(t *testing.T) {
	role1 := createRandomRole(t)
	err := testQueries.HardDeleteRole(context.Background(), role1.ID)
	require.NoError(t, err)

	role2, err := testQueries.GetRoleByID(context.Background(), role1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, role2)
}

func TestListRoles(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomRole(t)
	}

	arg := db.ListRolesParams{
		Limit:  5,
		Offset: 5,
	}

	roles, err := testQueries.ListRoles(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, roles, 5)

	for _, role := range roles {
		require.NotEmpty(t, role)
	}
}

func TestRoleExistsByName(t *testing.T) {
	// exists
	role := createRandomRole(t)

	exists, err := testQueries.RoleExistsByName(
		context.Background(),
		role.Name,
	)

	require.NoError(t, err)
	require.True(t, exists)

	// doesn't exist
	exists, err = testQueries.RoleExistsByName(
		context.Background(),
		"nonexistent-role",
	)

	require.NoError(t, err)
	require.False(t, exists)
}
