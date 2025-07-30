package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
	"userDataTransformer/internal/config"
)

type MiddlewareService struct {
	config *config.Config
	//database db.IDatabase
}

func NewMiddlewareService(config *config.Config) IMiddlewareService {
	return &MiddlewareService{config: config}
}

type IMiddlewareService interface {
	MiddlewareJWT() gin.HandlerFunc
	CreateServiceToken() (string, error)
	GetJWTSecret() []byte
}

func (ms *MiddlewareService) MiddlewareJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return ms.GetJWTSecret(), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
		if role, ok := claims["role"].(string); !ok || role != "internal-service" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		ctx.Set("serviceName", claims["iss"])
		ctx.Next()
	}
}

func (ms *MiddlewareService) CreateServiceToken() (string, error) {
	claims := jwt.MapClaims{
		"iss":  "user-data-transformer",
		"role": "internal-service",
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ms.GetJWTSecret())
}
func (ms *MiddlewareService) GetJWTSecret() []byte {
	return []byte(ms.config.JWT.Secret)
}
