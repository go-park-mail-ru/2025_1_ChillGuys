package csat

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
)

type ICsatUsecase interface {
	GetSurveyWithQuestions(ctx context.Context, name string) (*dto.SurveyWithQuestionsResponse, error)
	SubmitAnswer(ctx context.Context, req *dto.SubmitAnswersRequest) error
}