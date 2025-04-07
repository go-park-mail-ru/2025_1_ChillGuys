package cookie

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"net/http"
	"time"
)

// CookieProvider управляет настройками для работы с cookies
type CookieProvider struct {
	cfg *config.Config
}

// NewCookieProvider создает новый CookieProvider с конфигурацией
func NewCookieProvider(cfg *config.Config) *CookieProvider {
	return &CookieProvider{cfg: cfg}
}

// Set устанавливает cookie с заданным именем и значением токена
func (cp *CookieProvider) Set(w http.ResponseWriter, token, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Expires:  time.Now().UTC().Add(cp.cfg.JWTConfig.TokenLifeSpan),
	})
}

// Unset инвалидирует cookie с заданным именем
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
