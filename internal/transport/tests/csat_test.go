package tests

import (
	csat2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/csat/http"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	generatedcsat "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/csat"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/csat/mocks"
)

func TestGetAllSurveys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockSurveyServiceClient(ctrl)
	handler := csat2.NewCsatHandler(mockClient) // замените на фактический конструктор

	mockClient.EXPECT().
		GetAllSurveys(gomock.Any(), gomock.Any()).
		Return(&generatedcsat.SurveysList{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/csat/surveys", nil)
	w := httptest.NewRecorder()

	handler.GetAllSurveys(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestGetSurvey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockSurveyServiceClient(ctrl)
	handler := csat2.NewCsatHandler(mockClient)

	surveyName := "test-survey"
	expectedResponse := &generatedcsat.SurveyWithQuestionsResponse{
		SurveyId:    "11111111-1111-1111-1111-111111111111",
		Title:       "Test Survey",
		Description: "Test Desc",
		Questions: []*generatedcsat.QuestionResponseDTO{
			{QuestionId: "22222222-2222-2222-2222-222222222222", Text: "How are you?"},
		},
	}

	mockClient.EXPECT().
		GetSurveyWithQuestions(gomock.Any(), &generatedcsat.GetSurveyRequest{Name: surveyName}).
		Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/csat/survey/"+surveyName, nil)
	req = mux.SetURLVars(req, map[string]string{"name": surveyName})
	w := httptest.NewRecorder()

	handler.GetSurvey(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestGetSurveyStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockSurveyServiceClient(ctrl)
	handler := csat2.NewCsatHandler(mockClient)

	surveyID := "11111111-1111-1111-1111-111111111111"
	expectedResp := &generatedcsat.SurveyStatisticsResponse{
		Description: "Stats description",
		Questions: []*generatedcsat.QuestionStatisticsDTO{
			{
				QuestionId: "22222222-2222-2222-2222-222222222222",
				Text:       "How are you?",
				Stats:      []uint32{1, 2, 3},
			},
		},
	}

	mockClient.EXPECT().
		GetSurveyStatistics(gomock.Any(), &generatedcsat.GetStatisticsRequest{SurveyId: surveyID}).
		Return(expectedResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/csat/statistics/"+surveyID, nil)
	req = mux.SetURLVars(req, map[string]string{"surveyId": surveyID})
	w := httptest.NewRecorder()

	handler.GetSurveyStatistics(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
