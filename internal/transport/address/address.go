package address

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/validator"
	"github.com/mailru/easyjson"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/address"
	"github.com/google/uuid"
)

type AddressHandler struct {
	addressService address.IAddressUsecase
	geoapifyAPIKey string
	httpClient     *http.Client
}

func NewAddressHandler(
	u address.IAddressUsecase,
	geoapifyAPIKey string,
) *AddressHandler {
	return &AddressHandler{
		addressService: u,
		geoapifyAPIKey: geoapifyAPIKey,
		httpClient:     &http.Client{},
	}
}

// CreateAddress godoc
//
//	@Summary		Создание нового адреса
//	@Description	Создает новый адрес для текущего пользователя
//	@Tags			address
//	@Accept			json
//	@Produce		json
//	@Param			address			body	dto.AddressReqDTO	true	"Данные адреса"
//	@Param			X-Csrf-Token	header	string				true	"CSRF-токен для защиты от подделки запросов"
//	@Success		201				"Адрес успешно создан"
//	@Failure		400				{object}	object	"Неверный формат данных или ID пользователя"
//	@Failure		401				{object}	object	"Пользователь не авторизован"
//	@Failure		500				{object}	object	"Ошибка сервера при создании адреса"
//	@Security		TokenAuth
//	@Router			/addresses [post]
func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	const op = "AddressHandler.CreateAddress"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var createAddressReq dto.AddressDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &createAddressReq); err != nil {
		logger.WithError(err).Error("parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	if createAddressReq.Label.Valid && strings.TrimSpace(createAddressReq.Label.String) != "" {
		if err := validator.ValidateLabel(createAddressReq.Label.String); err != nil {
			logger.WithError(err).Error("invalid address label")
			response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
			return
		}
	}

	if !createAddressReq.AddressString.Valid || createAddressReq.AddressString.String == "" {
		logger.Error("address string is required")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "address string is required")
		return
	}

	userIDStr, ok := r.Context().Value(domains.UserIDKey{}).(string)
	if !ok {
		logger.Error("auth not found in context")
		response.SendJSONError(r.Context(), w, http.StatusUnauthorized, "auth not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).Error("parse user ID")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid auth id format")
		return
	}

	logger = logger.WithField("user_id", userID)

	//var geoData *dto.GeoapifyResponse
	//geoData, err = h.geocodeAddress(r.Context(), createAddressReq)
	//if err != nil {
	//	logger.WithError(err).Error("geoapify API error")
	//	response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to validate address")
	//	return
	//}
	//
	//var bestMatch *dto.GeoapifyFeature
	//for _, feature := range geoData.Features {
	//	if feature.Properties.ResultType == "building" && feature.Properties.Rank.Importance > 0.2 {
	//		if bestMatch == nil || feature.Properties.Rank.Importance > bestMatch.Properties.Rank.Importance {
	//			bestMatch = &feature
	//			break
	//		}
	//	}
	//}
	//
	//if bestMatch == nil {
	//	logger.Warn("no valid building address found")
	//	response.SendJSONError(r.Context(), w, http.StatusBadRequest, "no valid building address found")
	//	return
	//}

	if err := h.addressService.CreateAddress(r.Context(), userID, createAddressReq); err != nil {
		logger.WithError(err).Error("create address")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusCreated, nil)
}

func (h *AddressHandler) geocodeAddress(ctx context.Context, address dto.AddressDTO) (*dto.GeoapifyResponse, error) {
	const op = "AddressHandler.geocodeAddress"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	if !address.AddressString.Valid || address.AddressString.String == "" {
		logger.Error("address string is empty")
		return nil, fmt.Errorf("%s: address string is empty", op)
	}

	encodedAddress := url.QueryEscape(address.AddressString.String)

	apiURL := fmt.Sprintf("https://api.geoapify.com/v1/geocode/search?text=%s&apiKey=%s",
		encodedAddress,
		h.geoapifyAPIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		logger.WithError(err).Error("create request")
		return nil, fmt.Errorf("%s: failed to create request: %w", op, err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		logger.WithError(err).Error("call Geoapify API")
		return nil, fmt.Errorf("%s: failed to call Geoapify API: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Geoapify API returned status: %d", resp.StatusCode)
	}

	var geoResponse dto.GeoapifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&geoResponse); err != nil {
		logger.WithError(err).Error("decode Geoapify response")
		return nil, fmt.Errorf("%s: failed to decode Geoapify response: %w", op, err)
	}

	return &geoResponse, nil
}

// GetAddress godoc
//
//	@Summary		Получение списка адресов пользователя
//	@Description	Возвращает все адреса текущего пользователя
//	@Tags			address
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Успешный запрос"
//	@Failure		400	{object}	object				"Неверный формат ID пользователя"
//	@Failure		401	{object}	object				"Пользователь не авторизован"
//	@Failure		500	{object}	object				"Ошибка сервера при получении адресов"
//	@Security		TokenAuth
//	@Router			/addresses [get]
func (h *AddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	const op = "AddressHandler.GetAddress"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	userIDStr, isExist := r.Context().Value(domains.UserIDKey{}).(string)
	if !isExist {
		logger.Error("auth not found in context")
		response.SendJSONError(r.Context(), w, http.StatusUnauthorized, "auth not found in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).Error("parse user ID")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid auth id format")
		return
	}

	addresses, err := h.addressService.GetAddresses(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Error("get addresses")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	if addresses == nil {
		addresses = []dto.GetAddressResDTO{}
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string][]dto.GetAddressResDTO{
		"addresses": addresses,
	})
}

// GetPickupPoints godoc
//
//	@Summary		Получение списка пунктов выдачи
//	@Description	Возвращает все доступные пункты выдачи
//	@Tags			address
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Успешный запрос"
//	@Failure		500	{object}	object				"Ошибка сервера при получении пунктов выдачи"
//	@Router			/addresses/pickup-points [get]
func (h *AddressHandler) GetPickupPoints(w http.ResponseWriter, r *http.Request) {
	const op = "AddressHandler.GetPickupPoints"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	points, err := h.addressService.GetPickupPoints(r.Context())
	if err != nil {
		logger.WithError(err).Error("get pickup points")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string][]dto.GetPointAddressResDTO{
		"pickupPoints": points,
	})
}
