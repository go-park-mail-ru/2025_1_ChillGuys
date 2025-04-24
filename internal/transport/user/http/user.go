package user

import (
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	userClient gen.UserServiceClient
	config     *config.Config
}

func NewUserHandler(
	userClient gen.UserServiceClient,
	cfg *config.Config,
) *UserHandler {
	return &UserHandler{
		userClient: userClient,
		config:     cfg,
	}
}

// @Summary		Получить информацию о себе
// @Description	Возвращает информацию о текущем авторизованном пользователе
// @Tags			users
// @Produce		json
// @Success		200	{object}	dto.UserDTO	"Информация о пользователе"
// @Failure		400	{string}	string		"Некорректный запрос"
// @Failure		401	{string}	string		"Пользователь не авторизован"
// @Failure		500	{string}	string		"Внутренняя ошибка сервера"
// @Router			/users/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.GetMe"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	res, err := h.userClient.GetMe(r.Context(), &emptypb.Empty{})
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	user, err := dto.ConvertGrpcToUserDTO(res)
	if err != nil {
		logger.WithError(err).Error("failed to convert user")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to process user data")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, user)
}

// @Summary		Загрузить аватар
// @Description	Загружает изображение профиля пользователя
// @Tags			users
// @Accept			multipart/form-data
// @Produce		json
// @Param			file			formData	file				true	"Файл изображения"
// @Param			X-Csrf-Token	header		string				true	"CSRF-токен для защиты от подделки запросов"
// @Success		200				{object}	map[string]string	"URL загруженного аватара"
// @Failure		400				{string}	string				"Ошибка загрузки или обработки формы"
// @Failure		500				{string}	string				"Внутренняя ошибка сервера"
// @Router			/users/avatar [post]
func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UploadAvatar"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	if err := r.ParseMultipartForm(h.config.ServerConfig.MaxMultipartMemory); err != nil {
		logger.WithError(err).Error("failed to parse form data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	file, _, err := r.FormFile(h.config.ServerConfig.AvatarKey)
	if err != nil {
		logger.WithError(err).Error("no file uploaded")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "no file uploaded")
		return
	}
	defer file.Close()

	stream, err := h.userClient.UploadAvatar(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to create upload stream")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to start upload")
		return
	}

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.WithError(err).Error("failed to read file chunk")
			response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to read file")
			return
		}

		if err := stream.Send(&gen.BytesValue{Value: buf[:n]}); err != nil {
			logger.WithError(err).Error("failed to send chunk")
			response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to upload file")
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		logger.WithError(err).Error("failed to close stream")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to complete upload")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string]string{
		"imageURL": res.ImageURL,
	})
}

// @Summary		Обновить профиль пользователя
// @Description	Обновляет основную информацию пользователя
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			body			body		dto.UpdateUserProfileRequestDTO	true	"Данные для обновления профиля"
// @Param			X-Csrf-Token	header		string							true	"CSRF-токен для защиты от подделки запросов"
// @Success		200				{string}	string							"Профиль успешно обновлён"
// @Failure		400				{string}	string							"Невалидные данные"
// @Failure		500				{string}	string							"Ошибка при обновлении профиля"
// @Security		TokenAuth
// @Router			/users/update-profile [post]
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UpdateUserProfile"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var updateReq dto.UpdateUserProfileRequestDTO
	if err := request.ParseData(r, &updateReq); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.userClient.UpdateUserProfile(r.Context(), updateReq.ConvertToGrpcUpdateProfileReq())
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// @Summary		Обновить email пользователя
// @Description	Обновляет email текущего пользователя
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			body			body		dto.UpdateUserEmailDTO	true	"Новый email"
// @Param			X-Csrf-Token	header		string					true	"CSRF-токен для защиты от подделки запросов"
// @Success		200				{string}	string					"Email успешно обновлён"
// @Failure		400				{string}	string					"Невалидные данные"
// @Failure		500				{string}	string					"Ошибка при обновлении email"
// @Security		TokenAuth
// @Router			/users/update-email [post]
func (h *UserHandler) UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UpdateUserEmail"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var updateReq dto.UpdateUserEmailDTO
	if err := request.ParseData(r, &updateReq); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	_, err := h.userClient.UpdateUserEmail(r.Context(), &gen.UpdateUserEmailRequest{
		Email:    updateReq.Email,
		Password: updateReq.Password,
	})
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// @Summary		Обновить пароль пользователя
// @Description	Меняет пароль текущего пользователя
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			body			body		dto.UpdateUserPasswordDTO	true	"Старый и новый пароли"
// @Param			X-Csrf-Token	header		string						true	"CSRF-токен для защиты от подделки запросов"
// @Success		200				{string}	string						"Пароль успешно обновлён"
// @Failure		400				{string}	string						"Невалидные данные"
// @Failure		500				{string}	string						"Ошибка при обновлении пароля"
// @Security		TokenAuth
// @Router			/users/update-password [post]
func (h *UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UpdateUserPassword"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var updateReq dto.UpdateUserPasswordDTO
	if err := request.ParseData(r, &updateReq); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	_, err := h.userClient.UpdateUserPassword(r.Context(), &gen.UpdateUserPasswordRequest{
		OldPassword: updateReq.OldPassword,
		NewPassword: updateReq.NewPassword,
	})
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}
