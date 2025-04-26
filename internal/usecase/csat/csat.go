package csat

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

type ICsatRepository interface {
	GetSurvey(ctx context.Context, topicName string) (*models.SurveyWithQuestions, error)
	AddSurveySubmission(ctx context.Context, surveyID uuid.UUID, answers []models.Answer, userID uuid.UUID) error
	GetStatistics(context.Context, uuid.UUID) (*models.GetStatisticsResponse, error)
}

type CsatUsecase struct {
	repo ICsatRepository
}

func NewCsatUsecase(repo ICsatRepository) *CsatUsecase {
	return &CsatUsecase{
		repo: repo,
	}
}

func (u *CsatUsecase) GetSurveyWithQuestions(ctx context.Context, name string) (*dto.SurveyWithQuestionsResponse, error) {
	const op = "CsatUsecase.GetSurveyWithQuestions"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if name == "" {
		logger.Error("empty survey name")
		return nil, errs.ErrReadRequestData
	}

	survey, err := u.repo.GetSurvey(ctx, name)
	if err != nil {
		logger.WithError(err).Error("failed to get survey from repository")
		return nil, err
	}

	if survey == nil {
		logger.Error("survey not found")
		return nil, errs.NewNotFoundError("survey not found")
	}

	return &dto.SurveyWithQuestionsResponse{
		ID:          survey.ID,
		Title:       survey.Title,
		Description: survey.Description,
		Questions:   convertQuestionsToDTO(survey.Questions),
	}, nil
}

func (u *CsatUsecase) SubmitAnswer(ctx context.Context, req *dto.SubmitAnswersRequest) error {
	const op = "CsatUsecase.SubmitAnswer"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userIDStr, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist || userIDStr == "" {
		logger.Warn("user ID not found in context")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).Error("invalid user ID format")
		return fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
	}

	if req == nil {
		logger.Error("nil request")
		return errs.ErrParseRequestData
	}

	if len(req.Answers) == 0 {
		logger.Error("empty answers list")
		return errs.NewBusinessLogicError("empty answers list")
	}

	answers := make([]models.Answer, 0, len(req.Answers))
	for _, ans := range req.Answers {
		if ans.QuestionID == uuid.Nil {
			logger.Error("empty question ID in answer")
			return errs.NewBusinessLogicError("empty question ID in answer")
		}
		if ans.Value < 0 || ans.Value > 10 {
			logger.Error("invalid answer value")
			return errs.NewBusinessLogicError("invalid answer value")
		}

		answers = append(answers, models.Answer{
			QuestionID: ans.QuestionID,
			Value:      ans.Value,
		})
	}

	err = u.repo.AddSurveySubmission(ctx, req.SurveyID, answers, userID)
	if err != nil {
		logger.WithError(err).Error("failed to store answers")
		return err
	}

	return nil
}

func (u *CsatUsecase) GetSurveyStatistics(ctx context.Context, surveyID uuid.UUID) (*dto.SurveyStatisticsResponse, error) {
	const op = "CsatUsecase.GetSurveyStatistics"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if surveyID == uuid.Nil {
		logger.Error("empty survey ID")
		return nil, errs.NewBusinessLogicError("empty survey ID")
	}

	stats, err := u.repo.GetStatistics(ctx, surveyID)
	if err != nil {
		logger.WithError(err).Error("failed to get survey statistics")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if stats == nil {
		logger.Error("survey statistics not found")
		return nil, errs.NewNotFoundError("survey statistics not found")
	}

	return &dto.SurveyStatisticsResponse{
		Description: stats.Description,
		Questions:   convertStatisticsQuestionsToDTO(stats.Questions),
	}, nil
}

func convertStatisticsQuestionsToDTO(questions []models.QuestionStatistics) []dto.QuestionStatisticsDTO {
	result := make([]dto.QuestionStatisticsDTO, 0, len(questions))

	for _, q := range questions {
		stats := make([]uint, 11)

		for _, answer := range q.Answers {
			if answer <= 10 {
				stats[answer]++
			}
		}

		result = append(result, dto.QuestionStatisticsDTO{
			ID:    q.ID,
			Text:  q.Text,
			Stats: stats,
		})
	}

	return result
}

func convertQuestionsToDTO(questions []models.Question) []dto.QuestionResponseDTO {
	result := make([]dto.QuestionResponseDTO, 0, len(questions))
	for _, q := range questions {
		result = append(result, dto.QuestionResponseDTO{
			ID:   q.ID,
			Text: q.Text,
		})
	}
	return result
}
