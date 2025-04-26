package dto

import "github.com/google/uuid"

// Запрос на отправку ответов
type SubmitAnswersRequest struct {
	SurveyID uuid.UUID          `json:"surveyId"`
	Answers  []AnswerRequestDTO `json:"answers"`
}

// DTO для ответа в запросе
type AnswerRequestDTO struct {
	QuestionID uuid.UUID `json:"questionId"`
	Value      uint      `json:"value"`
}
