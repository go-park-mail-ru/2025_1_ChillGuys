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

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string][]dto.GetAddressResDTO{
		"addresses": addresses,
	})
}

func (h *AddressHandler) GetPickupPoints(w http.ResponseWriter, r *http.Request) {
	points, err := h.addressService.GetPickupPoints(r.Context())
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to get pickup points")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string][]dto.GetPointAddressResDTO{
		"pickupPoints": points,
	})
}
