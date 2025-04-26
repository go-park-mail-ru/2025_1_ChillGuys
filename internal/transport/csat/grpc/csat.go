package csat

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/csat"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/csat"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CsatGRPCHandler struct {
	gen.UnimplementedSurveyServiceServer
	csatUseCase csat.ICsatUsecase
}

func NewCsatGRPCHandler(u csat.ICsatUsecase) *CsatGRPCHandler {
	return &CsatGRPCHandler{
		csatUseCase: u,
	}
}

func (h *CsatGRPCHandler) GetSurveyWithQuestions(ctx context.Context, req *gen.GetSurveyRequest) (*gen.SurveyWithQuestionsResponse, error) {
	const op = "CsatGRPCHandler.GetSurveyWithQuestions"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	survey, err := h.csatUseCase.GetSurveyWithQuestions(ctx, req.Name)
	if err != nil {
		logger.WithError(err).Error("failed to get survey")
		return nil, errs.MapErrorToGRPC(err)
	}

	return dto.ConvertSurveyToGrpc(survey), nil
}

func (h *CsatGRPCHandler) SubmitAnswer(ctx context.Context, req *gen.SubmitAnswerRequest) (*emptypb.Empty, error) {
	const op = "CsatGRPCHandler.SubmitAnswer"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	answReq, err := dto.ConvertGrpcToSubmitRequest(req)
	if err != nil {
		logger.WithError(err).Error("failed sumbit answers")
		return &emptypb.Empty{}, nil
	}

	err = h.csatUseCase.SubmitAnswer(ctx, answReq)
	if err != nil {
		logger.WithError(err).Error("failed to submit answer")
		return nil, errs.MapErrorToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *CsatGRPCHandler) GetSurveyStatistics(ctx context.Context, req *gen.GetStatisticsRequest) (*gen.SurveyStatisticsResponse, error) {
	const op = "CsatGRPCHandler.GetSurveyStatistics"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	surveyID, err := uuid.Parse(req.SurveyId)
	if err != nil {
		logger.WithError(err).Error("invalid survey ID format")
		return nil, status.Error(codes.InvalidArgument, "invalid survey ID")
	}

	stats, err := h.csatUseCase.GetSurveyStatistics(ctx, surveyID)
	if err != nil {
		logger.WithError(err).Error("failed to get survey statistics")
		return nil, errs.MapErrorToGRPC(err)
	}

	return dto.ConvertModelsToGrpcStatisticsResponse(stats), nil
}


func (h *CsatGRPCHandler) GetAllSurveys(ctx context.Context, _ *emptypb.Empty) (*gen.SurveysList, error) {
	const op = "CsatGRPCHandler.GetAllSurveys"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	surveys, err := h.csatUseCase.GetAllSurveys(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to get surveys")
		return nil, errs.MapErrorToGRPC(err)
	}

	return dto.ConvertSurveysListToGrpc(surveys), nil
}