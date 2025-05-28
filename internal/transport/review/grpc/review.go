package review

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/review"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/review"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewGRPCServer struct {
	gen.UnimplementedReviewServiceServer
	reviewUsecase review.IReviewUsecase
}

func NewReviewGRPCServer (u review.IReviewUsecase) *ReviewGRPCServer {
	return &ReviewGRPCServer{
		reviewUsecase: u,
	}
}

func (s *ReviewGRPCServer) AddReview (ctx context.Context, req *gen.AddReviewRequest) (*gen.EmptyResponse, error) {
	const op = "ReviewGRPCServer.AddReview"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		logger.WithError(err).Error("invalid prouct ID format")
		return &gen.EmptyResponse{}, fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
	}

	rating := req.Rating
	if 1 < rating && rating > 5 {
		logger.Error("invalid rating")
		return &gen.EmptyResponse{}, fmt.Errorf("%s: %w", op, errs.NewBusinessLogicError("invalid rating"))
	}

	request := dto.AddReviewRequest {
		ProductID: productID,
		Rating: int(rating),
		Comment: req.Comment,
	}

	err = s.reviewUsecase.Add(ctx, request)
	if err != nil {
		logger.WithError(err).Error("add review")
		if errors.Is(err, errs.ErrAlreadyExists) {
			return &gen.EmptyResponse{}, status.Error(codes.AlreadyExists, "user has already reviewed this product")
		}
		return &gen.EmptyResponse{}, status.Error(codes.Internal, "internal server error")
	}

	return &gen.EmptyResponse{}, nil
}

func (s *ReviewGRPCServer) GetReviews (ctx context.Context, req *gen.GetReviewsRequest) (*gen.GetReviewsResponse, error) {
	const op = "ReviewGRPCServer.GetReviews"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		logger.WithError(err).Error("invalid product ID format")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
	}

	request := dto.GetReviewRequest {
		ProductID: productID,
		Offset: int(req.Offset),
	}

	reviews, err := s.reviewUsecase.Get(ctx, request)
	if err != nil {
		logger.WithError(err).Error("get reviews")
		return nil, fmt.Errorf("%s: %w", op, errs.NewBusinessLogicError("get reviews"))
	}

	return dto.ModelsToGRPC(reviews), nil
}