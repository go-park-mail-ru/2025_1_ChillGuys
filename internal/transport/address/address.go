package address

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/address"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

type GeoapifyResponse struct {
	Features []GeoapifyFeature `json:"features"`
}

type GeoapifyFeature struct {
	Properties struct {
		ResultType string  `json:"result_type"`
		Lon        float64 `json:"lon"`
		Lat        float64 `json:"lat"`
		Rank       struct {
			Importance float64 `json:"importance"`
		} `json:"rank"`
	} `json:"properties"`
}

type AddressHandler struct {
	addressService address.IAddressUsecase
	log            *logrus.Logger
	geoapifyAPIKey string
}

func NewAddressHandler(
	u address.IAddressUsecase,
	log *logrus.Logger,
	geoapifyAPIKey string,
) *AddressHandler {
	return &AddressHandler{
		addressService: u,
		log:            log,
		geoapifyAPIKey: geoapifyAPIKey,
	}
}
func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var createAddressReq dto.AddressDTO
	if err := request.ParseData(r, &createAddressReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	if !createAddressReq.City.Valid || createAddressReq.City.String == "" {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "city is required")
		return
	}
	if !createAddressReq.AddressString.Valid || createAddressReq.AddressString.String == "" {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "address string is required")
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

	geoData, err := h.geocodeAddress(r.Context(), createAddressReq)
	if err != nil {
		h.log.WithContext(r.Context()).Errorf("Geoapify API error: %v", err)
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to validate address")
		return
	}

	var bestMatch *GeoapifyFeature
	for _, feature := range geoData.Features {
		if feature.Properties.ResultType == "building" && feature.Properties.Rank.Importance > 0.2 {
			if bestMatch == nil || feature.Properties.Rank.Importance > bestMatch.Properties.Rank.Importance {
				bestMatch = &feature
			}
		}
	}

	if bestMatch == nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "no valid building address found")
		return
	}

	if err := h.addressService.CreateAddress(r.Context(), userID, createAddressReq); err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to create address")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusCreated, nil)
}

func (h *AddressHandler) geocodeAddress(ctx context.Context, address dto.AddressDTO) (*GeoapifyResponse, error) {
	addrParts := []string{}
	if address.AddressString.Valid {
		addrParts = append(addrParts, address.AddressString.String)
	}
	if address.City.Valid {
		addrParts = append(addrParts, address.City.String)
	}
	if address.Region.Valid {
		addrParts = append(addrParts, address.Region.String)
	}

	addressText := strings.Join(addrParts, ", ")
	encodedAddress := url.QueryEscape(addressText)

	apiURL := fmt.Sprintf("https://api.geoapify.com/v1/geocode/search?text=%s&apiKey=%s",
		encodedAddress,
		h.geoapifyAPIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Geoapify API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Geoapify API returned status: %d", resp.StatusCode)
	}

	var geoResponse GeoapifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Geoapify response: %w", err)
	}

	return &geoResponse, nil
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
