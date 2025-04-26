package models

import "github.com/google/uuid"

type Topic struct {
	ID uuid.UUID 
	Name string
}

type Survey struct {
	ID uuid.UUID
	TopicID uuid.UUID
	Title string
	Description string
}

type Question struct {
	ID uuid.UUID
	SurveyID uuid.UUID
	Text string
	Position uint
}

type Submssion struct {
	ID uuid.UUID
	UserID uuid.UUID
	SurveyID uuid.UUID
}

type Answer struct {
	id uuid.UUID
	SubmssionID uuid.UUID
	QuestionID uuid.UUID
	Value uint
}