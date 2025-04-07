package cookie

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"log"
	"net/http"
	"time"
)

type CookieProvider struct {
	cfg *config.Config
}

func NewCookieProvider(cfg *config.Config) *CookieProvider {
	if cfg == nil || cfg.JWTConfig == nil {
		log.Println("Warning: nil config or JWTConfig provided to CookieProvider")
	}
	return &CookieProvider{cfg: cfg}
}

func (cp *CookieProvider) Set(w http.ResponseWriter, token, name string) {
	if token == "" {
		log.Println("Warning: empty token for cookie", name)
		return
	}

	tokenLifeSpan := 24 * time.Hour
	if cp.cfg != nil && cp.cfg.JWTConfig != nil {
		tokenLifeSpan = cp.cfg.JWTConfig.TokenLifeSpan
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Expires:  time.Now().UTC().Add(tokenLifeSpan),
	})
}

func (cp *CookieProvider) Unset(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().UTC().AddDate(0, 0, -1),
		HttpOnly: true,
		Secure:   true,
	})
}
