package auth_test

import (
	"github.com/Nickolasll/urlshortener/internal/app/auth"
)

func Example() {
	token, err := auth.IssueToken()
	if err != nil {
		// Не удалось выпустить токен
	}
	isValid := auth.IsValid(token)
	if !isValid {
		// Только что выпущенный токен невалиден
	}
	userID := auth.GetUserID(token)
	if userID == "" {
		// Не удалось получить идентификатор пользователя
	}
}
