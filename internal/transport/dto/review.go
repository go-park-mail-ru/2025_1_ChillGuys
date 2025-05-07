package dto

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/review"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type AddReviewRequest struct {
	ProductID uuid.UUID `json:"productID" db:"product_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   string    `json:"comment" db:"comment"`
}

// ConvertAddReviewRequestToGRPC преобразует dto.AddReviewRequest в gen.AddReviewRequest
func ConvertAddReviewRequestToGRPC(dtoReq AddReviewRequest) *review.AddReviewRequest {
	return &review.AddReviewRequest{
		ProductId: dtoReq.ProductID.String(),
		Rating:    int32(dtoReq.Rating),
		Comment:   dtoReq.Comment,
	}
}

type GetReviewRequest struct {
	ProductID uuid.UUID `json:"productID" db:"review_id"`
	Offset    int       `json:"offset"`
}

// ConvertGetReviewRequestToGRPC преобразует dto.GetReviewRequest в gen.GetReviewsRequest
func ConvertGetReviewRequestToGRPC(dtoReq GetReviewRequest) *review.GetReviewsRequest {
	return &review.GetReviewsRequest{
		ProductId: dtoReq.ProductID.String(),
		Offset:    int32(dtoReq.Offset),
	}
}

type ReviewDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Surname  null.String    `json:"surname"`
	ImageURL null.String    `json:"imageURL"`
	Rating   int       `json:"rating"`
	Comment  string    `json:"comment"`
}

type ReviewsResponse struct {
	Reviews []ReviewDTO `json:"reviews"`
}

func ConvertToReviewDTO(review *models.Review) ReviewDTO {
	return ReviewDTO{
		ID:            review.ID,
		Name:          review.Name,
		Surname:       review.Surname,
		ImageURL:      review.ImageURL,
		Rating:        review.Rating,
		Comment:       review.Comment,
	}
}

func ConvertToReviewsResponse(reviews []*models.Review) ReviewsResponse {
	briefreviews := make([]ReviewDTO, 0, len(reviews))
	for _, review := range reviews {
		briefreviews = append(briefreviews, ConvertToReviewDTO(review))
	}

	return ReviewsResponse{
		Reviews: briefreviews,
	}
}

// ConvertReviewsResponseToProto преобразует dto.ReviewsResponse в gen.GetReviewsResponse
func ModelsToGRPC(dtoResp []*models.Review) *review.GetReviewsResponse {
	protoResp := &review.GetReviewsResponse{
		Reviews: make([]*review.Review, 0, len(dtoResp)),
	}

	for _, dtoReview := range dtoResp {
		var surname string
		if dtoReview.Surname.Valid {
			surname = dtoReview.Surname.String
		}

		var imageURL string
		if dtoReview.ImageURL.Valid {
			imageURL = dtoReview.ImageURL.String
		}

		protoReview := &review.Review{
			Id:       dtoReview.ID.String(),
			Name:     dtoReview.Name,
			Surname:  surname,
			ImageUrl: imageURL,
			Rating:   int32(dtoReview.Rating),
			Comment:  dtoReview.Comment,
		}
		protoResp.Reviews = append(protoResp.Reviews, protoReview)
	}

	return protoResp
}

// ConvertProtoToReviewsResponse преобразует gen.GetReviewsResponse в dto.ReviewsResponse
func ConvertGRPCToReviewsResponse(protoResp *review.GetReviewsResponse) ReviewsResponse {
	dtoResp := ReviewsResponse{
		Reviews: make([]ReviewDTO, 0, len(protoResp.GetReviews())),
	}

	for _, protoReview := range protoResp.GetReviews() {
		id, _ := uuid.Parse(protoReview.GetId())

		var surname null.String
		if protoReview.Surname != "" {
			surname = null.StringFrom(protoReview.Surname)
		}

		var imageURL null.String
		if protoReview.ImageUrl != "" {
			imageURL = null.StringFrom(protoReview.ImageUrl)
		}
		
		dtoReview := ReviewDTO{
			ID:       id,
			Name:     protoReview.GetName(),
			Surname:  surname,
			ImageURL: imageURL,
			Rating:   int(protoReview.GetRating()),
			Comment:  protoReview.GetComment(),
		}
		dtoResp.Reviews = append(dtoResp.Reviews, dtoReview)
	}

	return dtoResp
}