package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	csat "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/csat"
)

func TestSurveyRepository_GetSurvey(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := csat.NewSurveyRepository(db)

	t.Run("Success", func(t *testing.T) {
		surveyID := uuid.New()
		questionID := uuid.New()
		topicName := "test_topic"

		rows := sqlmock.NewRows([]string{
			"survey_id", "survey_title", "survey_description",
			"question_id", "question_position", "question_text",
		}).
			AddRow(
				surveyID, "Test Survey", "Survey Description",
				questionID, 1, "Question 1",
			).
			AddRow(
				surveyID, "Test Survey", "Survey Description",
				uuid.New(), 2, "Question 2",
			)

		mock.ExpectQuery("SELECT").
			WithArgs(topicName).
			WillReturnRows(rows)

		survey, err := repo.GetSurvey(context.Background(), topicName)
		assert.NoError(t, err)
		assert.Equal(t, surveyID, survey.ID)
		assert.Equal(t, "Test Survey", survey.Title)
		assert.Equal(t, "Survey Description", survey.Description)
		assert.Len(t, survey.Questions, 2)
		assert.Equal(t, questionID, survey.Questions[0].ID)
		assert.Equal(t, "Question 1", survey.Questions[0].Text)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WithArgs("unknown_topic").
			WillReturnRows(sqlmock.NewRows([]string{
				"survey_id", "survey_title", "survey_description",
				"question_id", "question_position", "question_text",
			}))

		_, err := repo.GetSurvey(context.Background(), "unknown_topic")
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QueryError", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WithArgs("test_topic").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetSurvey(context.Background(), "test_topic")
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"survey_id", "survey_title", "survey_description",
			"question_id", "question_position", "question_text",
		}).
			AddRow("invalid-uuid", nil, nil, nil, nil, nil)

		mock.ExpectQuery("SELECT").
			WithArgs("test_topic").
			WillReturnRows(rows)

		_, err := repo.GetSurvey(context.Background(), "test_topic")
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSurveyRepository_AddSurveySubmission(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    repo := csat.NewSurveyRepository(db)

    t.Run("Success", func(t *testing.T) {
        surveyID := uuid.New()
        userID := uuid.New()
        questionID := uuid.New()

        answers := []models.Answer{
            {QuestionID: questionID, Value: 5},
        }

        mock.ExpectBegin()
        mock.ExpectExec("INSERT INTO bazaar.submission").
            WithArgs(sqlmock.AnyArg(), userID, surveyID).
            WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectPrepare("INSERT INTO bazaar.answer").
            ExpectExec().
            WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), questionID, 5).
            WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectCommit()

        err := repo.AddSurveySubmission(context.Background(), surveyID, answers, userID)
        assert.NoError(t, err)

        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("TransactionError", func(t *testing.T) {
        mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

        err := repo.AddSurveySubmission(context.Background(), uuid.New(), []models.Answer{}, uuid.New())
        assert.Error(t, err)

        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("SubmissionInsertError", func(t *testing.T) {
        mock.ExpectBegin()
        mock.ExpectExec("INSERT INTO bazaar.submission").
            WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
            WillReturnError(sql.ErrConnDone)
        mock.ExpectRollback()

        err := repo.AddSurveySubmission(context.Background(), uuid.New(), []models.Answer{}, uuid.New())
        assert.Error(t, err)

        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("AnswerInsertError", func(t *testing.T) {
        mock.ExpectBegin()
        mock.ExpectExec("INSERT INTO bazaar.submission").
            WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
            WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectPrepare("INSERT INTO bazaar.answer").
            ExpectExec().
            WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
            WillReturnError(sql.ErrConnDone)
        mock.ExpectRollback()

        err := repo.AddSurveySubmission(context.Background(), uuid.New(), []models.Answer{{}}, uuid.New())
        assert.Error(t, err)

        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("PrepareStatementError", func(t *testing.T) {
        mock.ExpectBegin()
        mock.ExpectExec("INSERT INTO bazaar.submission").
            WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
            WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectPrepare("INSERT INTO bazaar.answer").
            WillReturnError(sql.ErrConnDone)
        mock.ExpectRollback()

        err := repo.AddSurveySubmission(context.Background(), uuid.New(), []models.Answer{{}}, uuid.New())
        assert.Error(t, err)

        assert.NoError(t, mock.ExpectationsWereMet())
    })
}

func TestSurveyRepository_GetStatistics(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    repo := csat.NewSurveyRepository(db)

    t.Run("Success", func(t *testing.T) {
        surveyID := uuid.New()
        questionID1 := uuid.New()
        questionID2 := uuid.New()

        rows := sqlmock.NewRows([]string{
            "survey_description", "question_id", "question_text", "answer_value",
        }).
            AddRow("Survey Desc", questionID1, "Question 1", 5).
            AddRow("Survey Desc", questionID1, "Question 1", 4).
            AddRow("Survey Desc", questionID2, "Question 2", sql.NullInt64{})

        mock.ExpectQuery("SELECT").
            WithArgs(surveyID).
            WillReturnRows(rows)

        stats, err := repo.GetStatistics(context.Background(), surveyID)
        assert.NoError(t, err)
        
        // Проверяем описание опроса
        assert.Equal(t, "Survey Desc", stats.Description)
        
        // Проверяем количество вопросов
        assert.Len(t, stats.Questions, 2)
        
        // Находим нужный вопрос по ID
        var question1, question2 *models.QuestionStatistics
        for i := range stats.Questions {
            if stats.Questions[i].ID == questionID1 {
                question1 = &stats.Questions[i]
            } else if stats.Questions[i].ID == questionID2 {
                question2 = &stats.Questions[i]
            }
        }
        
        // Проверяем первый вопрос
        assert.NotNil(t, question1, "Question 1 not found")
        if question1 != nil {
            assert.Equal(t, "Question 1", question1.Text)
            assert.Len(t, question1.Answers, 2)
            assert.Contains(t, question1.Answers, uint32(5))
            assert.Contains(t, question1.Answers, uint32(4))
        }
        
        // Проверяем второй вопрос
        assert.NotNil(t, question2, "Question 2 not found")
        if question2 != nil {
            assert.Equal(t, "Question 2", question2.Text)
            assert.Empty(t, question2.Answers)
        }

        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("EmptyResult", func(t *testing.T) {
        surveyID := uuid.New()
        
        rows := sqlmock.NewRows([]string{
            "survey_description", "question_id", "question_text", "answer_value",
        })

        mock.ExpectQuery("SELECT").
            WithArgs(surveyID).
            WillReturnRows(rows)

        stats, err := repo.GetStatistics(context.Background(), surveyID)
        assert.NoError(t, err)
        assert.Equal(t, "", stats.Description)
        assert.Empty(t, stats.Questions)

        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("QueryError", func(t *testing.T) {
        surveyID := uuid.New()

        mock.ExpectQuery("SELECT").
            WithArgs(surveyID).
            WillReturnError(sql.ErrConnDone)

        _, err := repo.GetStatistics(context.Background(), surveyID)
        assert.Error(t, err)

        assert.NoError(t, mock.ExpectationsWereMet())
    })
}

func TestSurveyRepository_GetAllSurvey(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := csat.NewSurveyRepository(db)

	t.Run("Success", func(t *testing.T) {
		surveyID1 := uuid.New()
		surveyID2 := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "title"}).
			AddRow(surveyID1, "Survey 1").
			AddRow(surveyID2, "Survey 2")

		mock.ExpectQuery("SELECT s.id, s.title FROM bazaar.survey s").
			WillReturnRows(rows)

		surveys, err := repo.GetAllSurvey(context.Background())
		assert.NoError(t, err)
		assert.Len(t, surveys, 2)
		assert.Equal(t, surveyID1, surveys[0].ID)
		assert.Equal(t, "Survey 1", surveys[0].Title)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title"})

		mock.ExpectQuery("SELECT s.id, s.title FROM bazaar.survey s").
			WillReturnRows(rows)

		surveys, err := repo.GetAllSurvey(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, surveys)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QueryError", func(t *testing.T) {
		mock.ExpectQuery("SELECT s.id, s.title FROM bazaar.survey s").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetAllSurvey(context.Background())
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title"}).
			AddRow("invalid-uuid", nil)

		mock.ExpectQuery("SELECT s.id, s.title FROM bazaar.survey s").
			WillReturnRows(rows)

		_, err := repo.GetAllSurvey(context.Background())
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}