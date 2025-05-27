package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/notification"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateNotification_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	n := models.Notification{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Text:   "Test notification",
		Title:  "Test",
		IsRead: false,
	}

	mock.ExpectExec("INSERT INTO bazaar.notification").
		WithArgs(n.ID, n.UserID, n.Text, n.Title, n.IsRead).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := notification.NewNotificationRepository(db)
	err = repo.Create(context.Background(), n)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNotification_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	n := models.Notification{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Text:   "Test notification",
		Title:  "Test",
		IsRead: false,
	}

	mock.ExpectExec("INSERT INTO bazaar.notification").
		WithArgs(n.ID, n.UserID, n.Text, n.Title, n.IsRead).
		WillReturnError(errors.New("database error"))

	repo := notification.NewNotificationRepository(db)
	err = repo.Create(context.Background(), n)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllByUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()
	now := time.Now()
	expected := []models.Notification{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Text:      "Test 1",
			Title:     "Test",
			IsRead:    false,
			UpdatedAt: now,
		},
		{
			ID:        uuid.New(),
			UserID:    userID,
			Text:      "Test 2",
			Title:     "Test",
			IsRead:    true,
			UpdatedAt: now.Add(-time.Hour),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "text", "title", "is_read", "updated_at"}).
		AddRow(expected[0].ID, expected[0].UserID, expected[0].Text, expected[0].Title, expected[0].IsRead, expected[0].UpdatedAt).
		AddRow(expected[1].ID, expected[1].UserID, expected[1].Text, expected[1].Title, expected[1].IsRead, expected[1].UpdatedAt)

	mock.ExpectQuery("SELECT id, user_id, text, title, is_read, updated_at FROM bazaar.notification WHERE user_id = \\$1 ORDER BY updated_at DESC LIMIT 10 OFFSET \\$2").
		WithArgs(userID, 0).
		WillReturnRows(rows)

	repo := notification.NewNotificationRepository(db)
	result, err := repo.GetAllByUser(context.Background(), userID, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expected[0].UserID, result[0].UserID)
	assert.Equal(t, expected[0].Text, result[0].Text)
	assert.Equal(t, expected[1].UserID, result[1].UserID)
	assert.Equal(t, expected[1].Text, result[1].Text)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllByUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()

	mock.ExpectQuery("SELECT id, user_id, text, title, is_read, updated_at FROM bazaar.notification WHERE user_id = \\$1 ORDER BY updated_at DESC LIMIT 10 OFFSET \\$2").
		WithArgs(userID, 0).
		WillReturnError(errors.New("database error"))

	repo := notification.NewNotificationRepository(db)
	_, err = repo.GetAllByUser(context.Background(), userID, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUnreadCount_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()
	expectedCount := 3

	rows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM bazaar.notification WHERE user_id = \\$1 AND is_read = false").
		WithArgs(userID).
		WillReturnRows(rows)

	repo := notification.NewNotificationRepository(db)
	count, err := repo.GetUnreadCount(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUnreadCount_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM bazaar.notification WHERE user_id = \\$1 AND is_read = false").
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	repo := notification.NewNotificationRepository(db)
	_, err = repo.GetUnreadCount(context.Background(), userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateReadStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	notificationID := uuid.New()
	isRead := true

	mock.ExpectExec("UPDATE bazaar.notification SET is_read = \\$1 WHERE id = \\$2").
		WithArgs(isRead, notificationID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := notification.NewNotificationRepository(db)
	err = repo.UpdateReadStatus(context.Background(), notificationID, isRead)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateReadStatus_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	notificationID := uuid.New()
	isRead := true

	mock.ExpectExec("UPDATE bazaar.notification SET is_read = \\$1 WHERE id = \\$2").
		WithArgs(isRead, notificationID).
		WillReturnError(errors.New("database error"))

	repo := notification.NewNotificationRepository(db)
	err = repo.UpdateReadStatus(context.Background(), notificationID, isRead)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
