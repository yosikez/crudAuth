package auth

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yosikez/crudAuth/database"
	"github.com/yosikez/crudAuth/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
	refreshTokenIssuer   = "crud-auth"
	refreshTokenAudience = "user"
)

var secretKey = os.Getenv("SECRET_KEY")

type Claims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateTokens(userId uint, username string) (accessToken, refreshToken string, err error) {
	accessTokenClaims := Claims{
		Id:       userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenDuration).Unix(),
		},
	}

	accessTokenToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err = accessTokenToken.SignedString([]byte(secretKey))

	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}
	userIDStr := strconv.Itoa(int(userId))

	refreshTokenClaims := Claims{
		Id:       userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenDuration).Unix(),
			Issuer:    refreshTokenIssuer,
			Audience:  refreshTokenAudience,
			Subject: userIDStr,
		},
	}

	refreshTokenToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshToken, err = refreshTokenToken.SignedString([]byte(secretKey))

	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func RefreshTokens(refreshToken string) (accessToken, newRefreshToken string, err error){

	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", errors.New("invalid refresh token signature")
		}

		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid || claims.Audience != refreshTokenAudience || claims.Issuer != refreshTokenIssuer {
		return "", "", errors.New("token is not valid")
	}

	accessTokenClaims := Claims{
		Id: claims.Id,
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenDuration).Unix(),
		},
	}

	accessTokenToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err = accessTokenToken.SignedString([]byte(secretKey))

	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}
	userIDStr := strconv.Itoa(int(claims.Id))

	newRefreshTokenClaims := Claims{
		Id: claims.Id,
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenDuration).Unix(),
			Issuer: refreshTokenIssuer,
			Audience: refreshTokenAudience,
			Subject: userIDStr,
		},
	}

	newRefreshTokenToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshTokenClaims)
	newRefreshToken, err = newRefreshTokenToken.SignedString([]byte(secretKey))

	if err != nil {
		return "", "", errors.New("failed to generate new refresh token")
	}

	return accessToken, newRefreshToken, nil	
}

func AuthenticateUser(username, password string) (*model.User, error){
	var user model.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		return nil, err
	}

	return &user, nil
}

func GetRefrehToken(userId uint) (*model.RefreshToken, error){
	var existingToken model.RefreshToken

	if err := database.DB.Where("user_id = ?", userId).First(&existingToken).Error; err != nil {
		return nil, err
	}

	return &existingToken, nil
}

func GetClaimsDataFromToken(c *gin.Context) (*Claims, error) {
    header := c.GetHeader("Authorization")
    if header == "" {
        return nil, errors.New("missing authorization header")
    }

    tokenStr := strings.Replace(header, "Bearer ", "", 1)

    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })

    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            return nil, errors.New("invalid token signature")
        }

        return nil, errors.New("invalid token")
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token claims")
    }

    return claims, nil
}


// func verifyToken(authHeader string) (*jwt.Token, error) {
// 	if !strings.HasPrefix(authHeader, "Bearer "){
// 		return nil, errors.New("invalid token format")
// 	}

// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	
// 	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(t *jwt.Token) (interface{}, error) {
// 		return []byte(secretKey), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return token, nil
// }
