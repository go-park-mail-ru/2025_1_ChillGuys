package address

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/address"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AddressHandler struct {
	addressService address.IAddressUsecase
	log            *logrus.Logger
}

func NewAddressHandler(
	u address.IAddressUsecase,
	log *logrus.Logger,
) *AddressHandler {
	return &AddressHandler{
		addressService: u,
		log:            log,
	}
}

// CreateAddress godoc
//
//	@Summary		Создание нового адреса
//	@Description	Создает новый адрес для текущего пользователя
//	@Tags			address
//	@Accept			json
//	@Produce		json
//	@Param			address	body	dto.AddressReqDTO	true	"Данные адреса"
//	@Success		201		"Адрес успешно создан"
//	@Failure		400		{object}	object	"Неверный формат данных или ID пользователя"
//	@Failure		401		{object}	object	"Пользователь не авторизован"
//	@Failure		500		{object}	object	"Ошибка сервера при создании адреса"
//	@Security		TokenAuth
//	@Router			/addresses [post]
func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var createAddressReq dto.AddressDTO
	if err := request.ParseData(r, &createAddressReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	userIDStr, ok := r.Context().Value(domains.UserIDKey{}).(string)
	if !ok {
		response.SendJSONError(r.Context(), w, http.StatusUnauthorized, "auth not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid auth id format")
		return
	}

	if err := h.addressService.CreateAddress(r.Context(), userID, createAddressReq); err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to create address")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusCreated, nil)
}

// GetAddress godoc
//
//	@Summary		Получение списка адресов пользователя
//	@Description	Возвращает все адреса текущего пользователя
//	@Tags			address
//	@Produce		json
//	@Success		200	{object}	dto.AddressListResponse	"Успешный запрос"
//	@Failure		400	{object}	object					"Неверный формат ID пользователя"
//	@Failure		401	{object}	object					"Пользователь не авторизован"
//	@Failure		500	{object}	object					"Ошибка сервера при получении адресов"
//	@Security		TokenAuth
//	@Router			/addresses [get]
func (h *AddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	userIDStr, isExist := r.Context().Value(domains.UserIDKey{}).(string)
	if !isExist {
		response.SendJSONError(r.Context(), w, http.StatusUnauthorized, "auth not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid auth id format")
		return
	}

	addresses, err := h.addressService.GetAddresses(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to get addresses")
		return
	}

	if addresses == nil {
		addresses = []dto.GetAddressResDTO{}
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, dto.AddressListResponse{
		Addresses: addresses,
	})
}

// GetPickupPoints godoc
//
//	@Summary		Получение списка пунктов выдачи
//	@Description	Возвращает все доступные пункты выдачи
//	@Tags			address
//	@Produce		json
//	@Success		200	{object}	dto.PickupPointListResponse	"Успешный запрос"
//	@Failure		500	{object}	object						"Ошибка сервера при получении пунктов выдачи"
//	@Router			/addresses/pickup-points [get]
func (h *AddressHandler) GetPickupPoints(w http.ResponseWriter, r *http.Request) {
	points, err := h.addressService.GetPickupPoints(r.Context())
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to get pickup points")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, dto.PickupPointListResponse{
		PickupPoints: points,
	})
}
