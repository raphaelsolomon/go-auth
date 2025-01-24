package middleware

import (
	"golang/internal/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	UserUuid string `json:"user_uuid"`
	jwt.RegisteredClaims
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authheader := ctx.GetHeader("Authorization")
		if authheader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authheader, "Bearer ")
		if tokenString == authheader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			ctx.Abort()
			return
		}

		secretkey := config.GetEnv("SECRET_KEY")

		token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretkey), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid/Expired token"})
			ctx.Abort()
			return
		}

		if claim, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			ctx.Set("user_uuid", claim.UserUuid)
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
