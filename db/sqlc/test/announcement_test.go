package test

import (
	"context"
	"database/sql"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAnnouncement(t *testing.T) db.Announcement {
	user := createRandomUser(t)

	arg := db.CreateAnnouncementParams{
		Name: utils.RandomString(utils.RandomInt(1, 100)),
		Description: sql.NullString{
			String: utils.RandomString(utils.RandomInt(1, 200)),
			Valid:  true,
		},
		ExternalLink: sql.NullString{Valid: false},
		CreatedBy:    user.ID,
	}

	announcement, err := testQueries.CreateAnnouncement(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, announcement)

	require.Equal(t, arg.Name, announcement.Name)
	require.Equal(t, arg.Description, announcement.Description)
	require.False(t, announcement.ExternalLink.Valid)
	require.Equal(t, arg.CreatedBy, announcement.CreatedBy)
	require.NotZero(t, announcement.ID)
	require.NotZero(t, announcement.CreatedAt)
	require.False(t, announcement.DeletedAt.Valid)

	return announcement
}

func TestCreateAnnouncement(t *testing.T) {
	createRandomAnnouncement(t)
}

func TestGetAnnouncementByID(t *testing.T) {
	announcement1 := createRandomAnnouncement(t)
	announcement2, err := testQueries.GetAnnouncementByID(context.Background(), announcement1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, announcement2)

	require.Equal(t, announcement1.ID, announcement2.ID)
	require.Equal(t, announcement1.Name, announcement2.Name)
	require.Equal(t, announcement1.Description, announcement2.Description)
	require.Equal(t, announcement1.ExternalLink, announcement2.ExternalLink)
	require.Equal(t, announcement1.CreatedBy, announcement2.CreatedBy)
	require.Equal(t, announcement1.CreatedAt, announcement2.CreatedAt)
}

func TestUpdateAnnouncement(t *testing.T) {
	announcement1 := createRandomAnnouncement(t)

	arg := db.UpdateAnnouncementParams{
		ID:   announcement1.ID,
		Name: utils.RandomString(utils.RandomInt(1, 100)),
		Description: sql.NullString{
			String: utils.RandomString(utils.RandomInt(1, 200)),
			Valid:  true,
		},
		ExternalLink: sql.NullString{
			String: "https://" + utils.RandomString(10) + ".com",
			Valid:  true,
		},
		CreatedBy: announcement1.CreatedBy,
	}

	announcement2, err := testQueries.UpdateAnnouncement(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, announcement2)

	require.Equal(t, announcement1.ID, announcement2.ID)
	require.Equal(t, arg.Name, announcement2.Name)
	require.Equal(t, arg.Description, announcement2.Description)
	require.Equal(t, arg.ExternalLink, announcement2.ExternalLink)
	require.True(t, announcement2.UpdatedAt.After(announcement1.UpdatedAt))
}

func TestDeleteAnnouncement(t *testing.T) {
	announcement1 := createRandomAnnouncement(t)
	err := testQueries.DeleteAnnouncement(context.Background(), announcement1.ID)
	require.NoError(t, err)

	announcement2, err := testQueries.GetAnnouncementByID(context.Background(), announcement1.ID)
	require.NoError(t, err)
	require.True(t, announcement2.DeletedAt.Valid)
}

func TestRestoreAnnouncement(t *testing.T) {
	announcement1 := createRandomAnnouncement(t)

	err := testQueries.DeleteAnnouncement(context.Background(), announcement1.ID)
	require.NoError(t, err)

	err = testQueries.RestoreAnnouncement(context.Background(), announcement1.ID)
	require.NoError(t, err)

	announcement2, err := testQueries.GetAnnouncementByID(context.Background(), announcement1.ID)
	require.NoError(t, err)
	require.False(t, announcement2.DeletedAt.Valid)
}

func TestHardDeleteAnnouncement(t *testing.T) {
	announcement1 := createRandomAnnouncement(t)
	err := testQueries.HardDeleteAnnouncement(context.Background(), announcement1.ID)
	require.NoError(t, err)

	announcement2, err := testQueries.GetAnnouncementByID(context.Background(), announcement1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, announcement2)
}

func TestListAnnouncements(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAnnouncement(t)
	}

	arg := db.ListAnnouncementsParams{
		Limit:  5,
		Offset: 5,
	}

	announcements, err := testQueries.ListAnnouncements(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, announcements, 5)

	for _, announcement := range announcements {
		require.NotEmpty(t, announcement)
		require.False(t, announcement.DeletedAt.Valid)
	}
}
