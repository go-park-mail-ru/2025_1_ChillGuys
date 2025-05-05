package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	review "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/review"
)

func TestReviewRepository_AddReview(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := review.NewReviewRepository(db)

	t.Run("Success", func(t *testing.T) {
		reviewDB := models.ReviewDB{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			ProductID: uuid.New(),
			Rating:    5,
			Comment:   "Great product!",
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE bazaar.product SET reviews_count = reviews_count + 1 WHERE id = $1").
			WithArgs(reviewDB.ProductID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO bazaar.review (id, user_id, product_id, rating, comment) VALUES ($1, $2, $3, $4, $5)").
			WithArgs(reviewDB.ID, reviewDB.UserID, reviewDB.ProductID, reviewDB.Rating, reviewDB.Comment).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.AddReview(context.Background(), reviewDB)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateCountError", func(t *testing.T) {
		reviewDB := models.ReviewDB{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			ProductID: uuid.New(),
			Rating:    4,
			Comment:   "Good product",
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE bazaar.product SET reviews_count = reviews_count + 1 WHERE id = $1").
			WithArgs(reviewDB.ProductID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.AddReview(context.Background(), reviewDB)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("InsertReviewError", func(t *testing.T) {
		reviewDB := models.ReviewDB{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			ProductID: uuid.New(),
			Rating:    3,
			Comment:   "Average product",
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE bazaar.product SET reviews_count = reviews_count + 1 WHERE id = $1").
			WithArgs(reviewDB.ProductID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO bazaar.review (id, user_id, product_id, rating, comment) VALUES ($1, $2, $3, $4, $5)").
			WithArgs(reviewDB.ID, reviewDB.UserID, reviewDB.ProductID, reviewDB.Rating, reviewDB.Comment).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.AddReview(context.Background(), reviewDB)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CommitError", func(t *testing.T) {
		reviewDB := models.ReviewDB{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			ProductID: uuid.New(),
			Rating:    2,
			Comment:   "Bad product",
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE bazaar.product SET reviews_count = reviews_count + 1 WHERE id = $1").
			WithArgs(reviewDB.ProductID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO bazaar.review (id, user_id, product_id, rating, comment) VALUES ($1, $2, $3, $4, $5)").
			WithArgs(reviewDB.ID, reviewDB.UserID, reviewDB.ProductID, reviewDB.Rating, reviewDB.Comment).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

		err := repo.AddReview(context.Background(), reviewDB)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReviewRepository_GetReview(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := review.NewReviewRepository(db)

	t.Run("Success", func(t *testing.T) {
		productID := uuid.New()
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "name", "surname", "image_url", "rating", "comment",
		}).
			AddRow(uuid.New(), "John", "Doe", "image1.jpg", 5, "Excellent").
			AddRow(uuid.New(), "Jane", "Smith", "image2.jpg", 4, "Good")

		expectedQuery := `SELECT 
			r.id, u.name, u.surname, u.image_url, r.rating, r.comment
		FROM bazaar.review r
		JOIN bazaar.user u ON r.user_id = u.id
		WHERE r.product_id = $1
		ORDER BY r.created_at DESC
        LIMIT 7 OFFSET $2`

		mock.ExpectQuery(expectedQuery).
			WithArgs(productID, offset).
			WillReturnRows(rows)

		reviews, err := repo.GetReview(context.Background(), productID, offset)
		assert.NoError(t, err)
		assert.Len(t, reviews, 2)
		assert.Equal(t, "John", reviews[0].Name)
		assert.Equal(t, "Jane", reviews[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NoReviews", func(t *testing.T) {
		productID := uuid.New()
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "name", "surname", "image_url", "rating", "comment",
		})

		expectedQuery := `SELECT 
			r.id, u.name, u.surname, u.image_url, r.rating, r.comment
		FROM bazaar.review r
		JOIN bazaar.user u ON r.user_id = u.id
		WHERE r.product_id = $1
		ORDER BY r.created_at DESC
        LIMIT 7 OFFSET $2`

		mock.ExpectQuery(expectedQuery).
			WithArgs(productID, offset).
			WillReturnRows(rows)

		reviews, err := repo.GetReview(context.Background(), productID, offset)
		assert.NoError(t, err)
		assert.Empty(t, reviews)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QueryError", func(t *testing.T) {
		productID := uuid.New()
		offset := 0

		expectedQuery := `SELECT 
			r.id, u.name, u.surname, u.image_url, r.rating, r.comment
		FROM bazaar.review r
		JOIN bazaar.user u ON r.user_id = u.id
		WHERE r.product_id = $1
		ORDER BY r.created_at DESC
        LIMIT 7 OFFSET $2`

		mock.ExpectQuery(expectedQuery).
			WithArgs(productID, offset).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetReview(context.Background(), productID, offset)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ScanError", func(t *testing.T) {
		productID := uuid.New()
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "name", "surname", "image_url", "rating", "comment",
		}).AddRow("invalid-uuid", nil, nil, nil, nil, nil)

		expectedQuery := `SELECT 
			r.id, u.name, u.surname, u.image_url, r.rating, r.comment
		FROM bazaar.review r
		JOIN bazaar.user u ON r.user_id = u.id
		WHERE r.product_id = $1
		ORDER BY r.created_at DESC
        LIMIT 7 OFFSET $2`

		mock.ExpectQuery(expectedQuery).
			WithArgs(productID, offset).
			WillReturnRows(rows)

		_, err := repo.GetReview(context.Background(), productID, offset)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}