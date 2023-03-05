package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yosikez/crudAuth/auth"
)
var secretKey = os.Getenv("SECRET_KEY")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			return
		}

		token, err := verifyToken(authHeader)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error" : "invalid token",
			})
			return
		}

		claims, ok :=  token.Claims.(*auth.Claims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error" : "invalid token",
			})
			return
		}

		c.Set("username", claims.Username)
		c.Set("userId", claims.Id)
		c.Next()
	}
}

func verifyToken(authHeader string) (*jwt.Token, error) {
	if !strings.HasPrefix(authHeader, "Bearer "){
		return nil, errors.New("invalid token format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}