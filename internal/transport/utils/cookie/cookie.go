package cookie

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
)

type CookieKeys string

const (
	Token     CookieKeys = "token"
	UserIDKey            = "userID"
)

func Cookie(w http.ResponseWriter, token, name string) {
	// Создаём cookie с переданными параметрами
	SSCookie := &http.Cookie{
		Name:     name,                                    // Имя куки
		Value:    token,                                   // Значение токена
		Path:     "/",                                     // Путь, для которого кука действительна
		SameSite: http.SameSiteStrictMode,                 // Политика SameSite
		HttpOnly: true,                                    // Доступность куки только для HTTP-запросов
		Expires:  time.Now().UTC().Add(jwt.TokenLifeSpan), // Время жизни куки
	}

	// Устанавливаем куку в ответ
	http.SetCookie(w, SSCookie)
}
