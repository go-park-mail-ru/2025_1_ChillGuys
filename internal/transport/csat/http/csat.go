package csat

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/csat"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CsatHandler struct {
	csatClient gen.SurveyServiceClient
}

func NewCsatHandler(csatClient gen.SurveyServiceClient) *CsatHandler {
	return &CsatHandler{
		csatClient: csatClient,
	}
}

func (h *CsatHandler) GetSurvey (w http.ResponseWriter, r *http.Request) {
	const op = "CsatHandler.Get"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	name := vars["name"]

	logger = logger.WithField("topic_name", name)

	res, err := h.csatClient.GetSurveyWithQuestions(r.Context(), &gen.GetSurveyRequest{Name:name})
	if err != nil {
		logger.WithError(err).Error("failed get survey")
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	survey, err := dto.ConvertGrpcToSurvey(res)
	if err != nil {
		logger.WithError(err).Error("failed to convert survey")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to process survey data")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, survey)
}

func (h *CsatHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	const op = "CsatHandler.SubmitAnswer"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var answerReq dto.SubmitAnswersRequest
	if err := request.ParseData(r, &answerReq); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.csatClient.SubmitAnswer(r.Context(), dto.ConvertToGrpcSubmitRequest(&answerReq))
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

func (h *CsatHandler) GetAllSurveys(w http.ResponseWriter, r *http.Request) {
	const op = "CsatHandler.GetAllSurveys"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	res, err := h.csatClient.GetAllSurveys(r.Context(), &emptypb.Empty{})
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	surveys, err := dto.ConvertGrpcToSurveyList(res)
	if err != nil {
		logger.WithError(err).Error("failed to convert surveys")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to process surveys data")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, surveys)
}