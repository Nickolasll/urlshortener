package auth

import (
	"time"

	"github.com/Nickolasll/urlshortener/internal/app/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// Выпускает новый токен для пользователя
func IssueToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.TokenExp) * time.Second)),
		},
		UserID: uuid.New().String(),
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Проверяет валидность токена
func IsValid(tokenString string) bool {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		},
	)
	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}
	return true
}

// Получает идентификатор пользователя из токена
func GetUserID(tokenString string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		},
	)
	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}
	return claims.UserID
}
