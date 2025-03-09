package utils

import (
	"net/http"
	"time"
)

type CookieKeys string

const (
	Token          CookieKeys = "token"
	UserIDKey                 = "userID"
	UserVersionKey            = "userVersion"
)

func Cookie(w http.ResponseWriter, token, name string) {
	// Создаём cookie с переданными параметрами
	SSCookie := &http.Cookie{
		Name:     name,                                 // Имя куки
		Value:    token,                                // Значение токена
		Path:     "/",                                  // Путь, для которого кука действительна
		SameSite: http.SameSiteStrictMode,              // Политика SameSite
		HttpOnly: true,                                 // Доступность куки только для HTTP-запросов
		Expires:  time.Now().UTC().Add(time.Hour * 24), // Время жизни куки (1 день)
	}

	// Устанавливаем куку в ответ
	http.SetCookie(w, SSCookie)
}
