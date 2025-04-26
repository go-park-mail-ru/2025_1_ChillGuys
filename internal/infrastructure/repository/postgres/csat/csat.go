package csat

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

const (
	queryGetSurvey = `SELECT
		s.id AS survey_id,
		s.title AS survey_title,
		s.description AS survey_description,
		q.id AS question_id,
		q.position AS question_position,
		q.text AS question_text
		FROM bazaar.topic t
		JOIN bazaar.survey s ON s.topic_id = t.id
		JOIN bazaar.question q ON q.survey_id = s.id
		WHERE t.name = $1
		ORDER BY q.position;`
	insertAnswerQuery = `
		INSERT INTO bazaar.answer (id, submission_id, question_id, value)
		VALUES ($1, $2, $3, $4);`
	insertSubmissionQuery = `
		INSERT INTO bazaar.submission (id, user_id, survey_id)
		VALUES ($1, $2, $3);`
)

type ISurveyRepository interface {
	GetSurvey(ctx context.Context, topicName string) (models.SurveyWithQuestions, error)
	AddSurveySubmission(ctx context.Context, answer dto.SubmitAnswersRequest, userID uuid.UUID) error
}

type SurveyRepository struct {
	db *sql.DB
}

func NewSurveyRepository(db *sql.DB) *SurveyRepository {
	return &SurveyRepository{
		db: db,
	}
}

func (r *SurveyRepository) GetSurvey(ctx context.Context, topicName string) (models.SurveyWithQuestions, error) {
	const op = "SurveyRepository.GetSurvey"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	rows, err := r.db.QueryContext(ctx, queryGetSurvey, topicName)
	if err != nil {
		logger.WithError(err).Error("failed to query survey")
		return models.SurveyWithQuestions{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var survey models.SurveyWithQuestions
	survey.Questions = make([]models.Question, 0)

	first := true

	for rows.Next() {
		var surveyData models.SurveyQuestionData

		if err = rows.Scan(
			&surveyData.SurveyID,
			&surveyData.SurveyTitle,
			&surveyData.SurveyDescription,
			&surveyData.QuestionID,
			&surveyData.QuestionPosition,
			&surveyData.QuestionText,
		); err != nil {
			logger.WithError(err).Error("failed to scan survey data")
			return models.SurveyWithQuestions{}, fmt.Errorf("scan survey: %w", err)
		}

		if first {
			survey.ID = surveyData.SurveyID
			survey.Title = surveyData.SurveyTitle
			survey.Description = surveyData.SurveyDescription
			first = false
		}

		survey.Questions = append(survey.Questions, models.Question{
			ID:       surveyData.QuestionID,
			Position: surveyData.QuestionPosition,
			Text:     surveyData.QuestionText,
		})
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return models.SurveyWithQuestions{}, fmt.Errorf("rows error: %w", err)
	}

	if len(survey.Questions) == 0 {
		logger.Warn("no questions found for survey")
		return models.SurveyWithQuestions{}, errs.NewNotFoundError("survey not found for topic")
	}

	return survey, nil
}

func (r *SurveyRepository) AddSurveySubmission(ctx context.Context, answer dto.SubmitAnswersRequest, userID uuid.UUID) error {
	const op = "SurveyRepository.AddSurveySubmission"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("failed to begin transaction")
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	submissionID := uuid.New()

	_, err = tx.ExecContext(ctx, insertSubmissionQuery, submissionID, userID, answer.SurveyID)
	if err != nil {
		logger.WithError(err).Error("failed to insert submission")
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.PrepareContext(ctx, insertAnswerQuery)
	if err != nil {
		logger.WithError(err).Error("failed to prepare insert answer statement")
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	for _, ans := range answer.Answers {
		answerID := uuid.New()
		_, err = stmt.ExecContext(ctx, answerID, submissionID, ans.QuestionID, ans.Value)
		if err != nil {
			logger.WithError(err).Error("failed to insert answer")
			tx.Rollback()
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("failed to commit transaction")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
