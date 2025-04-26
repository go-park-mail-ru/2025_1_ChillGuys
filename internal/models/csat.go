package models

import "github.com/google/uuid"

type Topic struct {
	ID   uuid.UUID
	Name string
}

type Survey struct {
	ID          uuid.UUID
	TopicID     uuid.UUID
	Title       string
	Description string
}

type Question struct {
	ID       uuid.UUID
	SurveyID uuid.UUID
	Text     string
	Position uint
}

type Submission struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	SurveyID uuid.UUID
}

type Answer struct {
	id          uuid.UUID
	SubmssionID uuid.UUID
	QuestionID  uuid.UUID
	Value       uint
}

type SurveyWithQuestions struct {
	ID          uuid.UUID
	Title       string
	Description string
	Questions   []Question
}

type SurveyQuestionData struct {
	SurveyID          uuid.UUID
	SurveyTitle       string
	SurveyDescription string
	QuestionID        uuid.UUID
	QuestionPosition  uint
	QuestionText      string
}

type GetStatisticsResponse struct {
	Description string
	Questions   []QuestionStatistics `json:"questions"`
}

type QuestionStatistics struct {
	ID      uuid.UUID `json:"questionId"`
	Text    string    `json:"text"`
	Answers []uint32  `json:"answer"`
}
