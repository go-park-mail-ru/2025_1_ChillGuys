package dto

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/csat"
	"github.com/google/uuid"
)

type SurveyWithQuestionsResponse struct {
	ID          uuid.UUID             `json:"surveyId"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Questions   []QuestionResponseDTO `json:"questions"`
}

func ConvertGrpcToSurvey(s *gen.SurveyWithQuestionsResponse) (*SurveyWithQuestionsResponse, error) {
	surveyID, err := uuid.Parse(s.SurveyId)
	if err != nil {
		return nil, err
	}

	questions := make([]QuestionResponseDTO, 0, len(s.Questions))
	for _, q := range s.Questions {
		questionID, err := uuid.Parse(q.QuestionId)
		if err != nil {
			return nil, err
		}

		questions = append(questions, QuestionResponseDTO{
			ID:   questionID,
			Text: q.Text,
		})
	}
	return &SurveyWithQuestionsResponse{
		ID:          surveyID,
		Title:       s.Title,
		Description: s.Description,
		Questions:   questions,
	}, nil
}

func ConvertSurveyToGrpc(s *SurveyWithQuestionsResponse) *gen.SurveyWithQuestionsResponse {
	questions := make([]*gen.QuestionResponseDTO, 0, len(s.Questions))
	for _, q := range s.Questions {
		questions = append(questions, &gen.QuestionResponseDTO{
			QuestionId: q.ID.String(),
			Text:       q.Text,
		})
	}

	return &gen.SurveyWithQuestionsResponse{
		SurveyId:    s.ID.String(),
		Title:       s.Title,
		Description: s.Description,
		Questions:   questions,
	}
}

type QuestionResponseDTO struct {
	ID   uuid.UUID `json:"questionId"`
	Text string    `json:"text"`
}

// Запрос на отправку ответов
type SubmitAnswersRequest struct {
	SurveyID  uuid.UUID          `json:"surveyId"`
	Answers   []AnswerRequestDTO `json:"answers"`
}

func ConvertToGrpcSubmitRequest(s *SubmitAnswersRequest) *gen.SubmitAnswerRequest {
	answers := make([]*gen.AnswerRequestDTO, 0, len(s.Answers))
	for _, a := range s.Answers {
		answers = append(answers, &gen.AnswerRequestDTO{
			QuestionId: a.QuestionID.String(),
			Value:      uint32(a.Value),
		})
	}

	return &gen.SubmitAnswerRequest{
		SurveyId: s.SurveyID.String(),
		Answers:  answers,
	}
}

// ConvertGrpcToSubmitRequest конвертирует gRPC SubmitAnswerRequest в DTO SubmitAnswersRequest
func ConvertGrpcToSubmitRequest(grpcReq *gen.SubmitAnswerRequest) (*SubmitAnswersRequest, error) {
    surveyID, err := uuid.Parse(grpcReq.SurveyId)
    if err != nil {
        return nil, err
    }

    answers := make([]AnswerRequestDTO, 0, len(grpcReq.Answers))
    for _, grpcAnswer := range grpcReq.Answers {
        questionID, err := uuid.Parse(grpcAnswer.QuestionId)
        if err != nil {
            return nil, err
        }

        answers = append(answers, AnswerRequestDTO{
            QuestionID: questionID,
            Value:      uint(grpcAnswer.Value),
        })
    }

    return &SubmitAnswersRequest{
        SurveyID: surveyID,
        Answers:  answers,
    }, nil
}

// DTO для ответа в запросе
type AnswerRequestDTO struct {
	QuestionID uuid.UUID `json:"questionId"`
	Value      uint      `json:"value"`
}

type GetStatisticsResponse struct {
	Questions []QuestionStatisticsDTO `json:"questions"`
}

type QuestionStatisticsDTO struct {
	ID      uuid.UUID `json:"questionId"`
	Text    string    `json:"text"`
	Answers []uint    `json:"answer"`
}

type BriefSurveyDTO struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type SurveysListDTO struct {
	Surveys []BriefSurveyDTO `json:"surveys"`
}

func ConvertModelsToSurveysListDTO(surveys []models.Survey) *SurveysListDTO {
	result := &SurveysListDTO{
		Surveys: make([]BriefSurveyDTO, 0, len(surveys)),
	}

	for _, s := range surveys {
		result.Surveys = append(result.Surveys, BriefSurveyDTO{
			ID:    s.ID,
			Title: s.Title,
		})
	}

	return result
}

func ConvertSurveysListToGrpc(s *SurveysListDTO) *gen.SurveysList {
	surveys := make([]*gen.BriefSurvey, 0, len(s.Surveys))
	for _, survey := range s.Surveys {
		surveys = append(surveys, &gen.BriefSurvey{
			Id:    survey.ID.String(),
			Title: survey.Title,
		})
	}

	return &gen.SurveysList{
		Surveys: surveys,
	}
}

func ConvertGrpcToSurveyList(grpcList *gen.SurveysList) (*SurveysListDTO, error) {
	surveys := make([]BriefSurveyDTO, 0, len(grpcList.Surveys))
	
	for _, grpcSurvey := range grpcList.Surveys {
		id, err := uuid.Parse(grpcSurvey.Id)
		if err != nil {
			return nil, err
		}

		surveys = append(surveys, BriefSurveyDTO{
			ID:    id,
			Title: grpcSurvey.Title,
		})
	}

	return &SurveysListDTO{
		Surveys: surveys,
	}, nil
}