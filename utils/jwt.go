package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Структура для JWT-токена
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT создает JWT-токен для пользователя
func GenerateJWT(userID uint, role string) (string, error) {
	// Получение секретного ключа из переменных окружения
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("jwt secret is not set")
	}

	// Установка срока действия токена (7 дней)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Создание claims
	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписание токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT разбирает и проверяет JWT-токен
func ParseJWT(tokenStr string) (uint, string, error) {
	// Получение секретного ключа из переменных окружения
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return 0, "", errors.New("jwt secret is not set")
	}

	// Парсинг токена
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	// Проверка валидности токена
	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	// Получение claims из токена
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return 0, "", errors.New("invalid token claims")
	}

	return claims.UserID, claims.Role, nil
}
