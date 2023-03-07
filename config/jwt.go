package config

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
	RefreshTokenIssuer   = "crud-auth"
	RefreshTokenAudience = "user"
)

var SecretKey = os.Getenv("SECRET_KEY")

type Claims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
