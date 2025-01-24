package jwt

import (
	"golang/internal/config"
	"golang/internal/middleware"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateSignToken(uuid string) (string, error) {
	claims := &middleware.TokenClaims{
		UserUuid: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	secretKey := config.GetEnv("SECRET_KEY")
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
}
