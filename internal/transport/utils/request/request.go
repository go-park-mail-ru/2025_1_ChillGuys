package request

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"io"
	"net/http"
)

func ParseData(r *http.Request, request any) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errs.ErrReadRequestData
	}

	if err := json.Unmarshal(body, request); err != nil {
		return errs.ErrParseRequestData
	}

	return nil
}
