package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	csat "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/csat"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetSurveyWithQuestions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	surveyName := "test-survey"

	expectedSurvey := &models.SurveyWithQuestions{
		ID:          uuid.New(),
		Title:       "Test Survey",
		Description: "Test Description",
		Questions: []models.Question{
			{
				ID:   uuid.New(),
				Text: "Question 1",
			},
			{
				ID:   uuid.New(),
				Text: "Question 2",
			},
		},
	}

	mockRepo.EXPECT().
		GetSurvey(ctx, surveyName).
		Return(expectedSurvey, nil)

	result, err := usecase.GetSurveyWithQuestions(ctx, surveyName)

	assert.NoError(t, err)
	assert.Equal(t, expectedSurvey.ID, result.ID)
	assert.Equal(t, expectedSurvey.Title, result.Title)
	assert.Equal(t, expectedSurvey.Description, result.Description)
	assert.Len(t, result.Questions, len(expectedSurvey.Questions))
}

func TestGetSurveyWithQuestions_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	_, err := usecase.GetSurveyWithQuestions(ctx, "")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrReadRequestData))
}

func TestGetSurveyWithQuestions_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	surveyName := "non-existent-survey"

	mockRepo.EXPECT().
		GetSurvey(ctx, surveyName).
		Return(nil, nil)

	_, err := usecase.GetSurveyWithQuestions(ctx, surveyName)

	assert.Error(t, err)
	assert.True(t, errors.Is(errs.ErrNotFound, errs.ErrNotFound))
}

func TestGetSurveyWithQuestions_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	surveyName := "test-survey"

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		GetSurvey(ctx, surveyName).
		Return(nil, expectedError)

	_, err := usecase.GetSurveyWithQuestions(ctx, surveyName)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestSubmitAnswer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	userID := uuid.New()
	ctx = context.WithValue(ctx, domains.UserIDKey{}, userID.String())

	req := &dto.SubmitAnswersRequest{
		SurveyID: uuid.New(),
		Answers: []dto.AnswerRequestDTO{
			{
				QuestionID: uuid.New(),
				Value:      5,
			},
			{
				QuestionID: uuid.New(),
				Value:      8,
			},
		},
	}

	mockRepo.EXPECT().
		AddSurveySubmission(ctx, req.SurveyID, gomock.Any(), userID).
		Return(nil)

	err := usecase.SubmitAnswer(ctx, req)

	assert.NoError(t, err)
}

func TestSubmitAnswer_UserIDNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	req := &dto.SubmitAnswersRequest{
		SurveyID: uuid.New(),
		Answers: []dto.AnswerRequestDTO{
			{
				QuestionID: uuid.New(),
				Value:      5,
			},
		},
	}

	err := usecase.SubmitAnswer(ctx, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
}

func TestSubmitAnswer_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	ctx = context.WithValue(ctx, domains.UserIDKey{}, "invalid-uuid")

	req := &dto.SubmitAnswersRequest{
		SurveyID: uuid.New(),
		Answers: []dto.AnswerRequestDTO{
			{
				QuestionID: uuid.New(),
				Value:      5,
			},
		},
	}

	err := usecase.SubmitAnswer(ctx, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidID))
}

func TestSubmitAnswer_NilRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	userID := uuid.New()
	ctx = context.WithValue(ctx, domains.UserIDKey{}, userID.String())

	err := usecase.SubmitAnswer(ctx, nil)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrParseRequestData))
}

func TestSubmitAnswer_EmptyAnswers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	userID := uuid.New()
	ctx = context.WithValue(ctx, domains.UserIDKey{}, userID.String())

	req := &dto.SubmitAnswersRequest{
		SurveyID: uuid.New(),
		Answers:  []dto.AnswerRequestDTO{},
	}

	err := usecase.SubmitAnswer(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty answers list")
}

func TestSubmitAnswer_InvalidQuestionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	userID := uuid.New()
	ctx = context.WithValue(ctx, domains.UserIDKey{}, userID.String())

	req := &dto.SubmitAnswersRequest{
		SurveyID: uuid.New(),
		Answers: []dto.AnswerRequestDTO{
			{
				QuestionID: uuid.Nil,
				Value:      5,
			},
		},
	}

	err := usecase.SubmitAnswer(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty question ID in answer")
}

func TestSubmitAnswer_InvalidAnswerValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	userID := uuid.New()
	ctx = context.WithValue(ctx, domains.UserIDKey{}, userID.String())

	req := &dto.SubmitAnswersRequest{
		SurveyID: uuid.New(),
		Answers: []dto.AnswerRequestDTO{
			{
				QuestionID: uuid.New(),
				Value:      11, // invalid value
			},
		},
	}

	err := usecase.SubmitAnswer(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid answer value")
}

func TestSubmitAnswer_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	userID := uuid.New()
	ctx = context.WithValue(ctx, domains.UserIDKey{}, userID.String())

	req := &dto.SubmitAnswersRequest{
		SurveyID: uuid.New(),
		Answers: []dto.AnswerRequestDTO{
			{
				QuestionID: uuid.New(),
				Value:      5,
			},
		},
	}

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		AddSurveySubmission(ctx, req.SurveyID, gomock.Any(), userID).
		Return(expectedError)

	err := usecase.SubmitAnswer(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestGetSurveyStatistics_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	surveyID := uuid.New()

	expectedStats := &models.GetStatisticsResponse{
		Description: "Test Survey Statistics",
		Questions: []models.QuestionStatistics{
			{
				ID:   uuid.New(),
				Text: "Question 1",
				Answers: []uint32{
					5, 7, 8, 9, 10,
				},
			},
		},
	}

	mockRepo.EXPECT().
		GetStatistics(ctx, surveyID).
		Return(expectedStats, nil)

	result, err := usecase.GetSurveyStatistics(ctx, surveyID)

	assert.NoError(t, err)
	assert.Equal(t, expectedStats.Description, result.Description)
	assert.Len(t, result.Questions, len(expectedStats.Questions))
}

func TestGetSurveyStatistics_EmptySurveyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	_, err := usecase.GetSurveyStatistics(ctx, uuid.Nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty survey ID")
}

func TestGetSurveyStatistics_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	surveyID := uuid.New()

	mockRepo.EXPECT().
		GetStatistics(ctx, surveyID).
		Return(nil, nil)

	_, err := usecase.GetSurveyStatistics(ctx, surveyID)

	assert.Error(t, err)
	assert.True(t, errors.Is(errs.ErrNotFound, errs.ErrNotFound))
}

func TestGetSurveyStatistics_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	surveyID := uuid.New()

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		GetStatistics(ctx, surveyID).
		Return(nil, expectedError)

	_, err := usecase.GetSurveyStatistics(ctx, surveyID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestGetAllSurveys_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	expectedSurveys := []models.Survey{
		{
			ID:          uuid.New(),
			Title:       "Survey 1",
			Description: "Description 1",
		},
		{
			ID:          uuid.New(),
			Title:       "Survey 2",
			Description: "Description 2",
		},
	}

	mockRepo.EXPECT().
		GetAllSurvey(ctx).
		Return(expectedSurveys, nil)

	result, err := usecase.GetAllSurveys(ctx)

	assert.NoError(t, err)
	assert.Len(t, result.Surveys, len(expectedSurveys))
}

func TestGetAllSurveys_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockICsatRepository(ctrl)
	usecase := csat.NewCsatUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		GetAllSurvey(ctx).
		Return(nil, expectedError)

	_, err := usecase.GetAllSurveys(ctx)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInternal))
}