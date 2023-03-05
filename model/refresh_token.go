package model

import (
	"os"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	Token    string `gorm:"refresh_token" json:"refresh_token"`
	UserId   uint   `gorm:"user_id"`
	Username string `gorm:"username"`
}

type TokenClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

var secretKey = os.Getenv("SECRET_KEY")

func (t *RefreshToken) IsValid(refreshToken string) bool {
	token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return false
	}

	return t.UserId == claims.Id
}
