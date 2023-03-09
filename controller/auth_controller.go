package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yosikez/crudAuth/auth"
	"github.com/yosikez/crudAuth/database"
	"github.com/yosikez/crudAuth/input"
	"github.com/yosikez/crudAuth/model"

	cusMessage "github.com/yosikez/custom-error-message"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (a *AuthController) Register(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		errFields := cusMessage.GetErrMess(err, user, nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation error",
			"errors":  errFields,
		})

		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create user",
			"error":   err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens(user.Id, user.Username, user.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed",
			"error":   err.Error(),
		})
		return
	}

	if err := database.DB.Create(&model.RefreshToken{Token: refreshToken, UserId: user.Id}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save refresh token",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (a *AuthController) Login(c *gin.Context) {
	var body input.LoginInput

	if err := c.ShouldBindJSON(&body); err != nil {
		errFields := cusMessage.GetErrMess(err, body, nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation error",
			"errors":  errFields,
		})
		return
	}

	user, err := auth.AuthenticateUser(body.Username, body.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to login",
			"error":   "invalid credentials",
		})

		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens(user.Id, user.Username, user.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to login",
			"error":   err.Error(),
		})

		return
	}

	existingToken, err := auth.GetRefrehToken(user.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get refresh token",
			"error":   err.Error(),
		})

		return
	}

	if existingToken != nil {
		existingToken.Token = refreshToken

		if err := database.DB.Save(&existingToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to update refresh token",
				"error":   err.Error(),
			})

			return
		}
	} else {
		newRefreshToken := model.RefreshToken{
			UserId:   user.Id,
			Token:    refreshToken,
			Username: user.Username,
		}

		if err := database.DB.Create(&newRefreshToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to save refresh token",
				"error":   err.Error(),
			})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (a *AuthController) RefreshToken(c *gin.Context) {
	var resfreshTokenRequest model.RefreshToken

	if err := c.ShouldBindJSON(&resfreshTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	claims, err := auth.GetClaimsDataFromToken(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	existingToken, err := auth.GetRefrehToken(claims.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if !existingToken.IsValid(resfreshTokenRequest.Token) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid refresh token",
			"test":  "test",
		})
		return
	}

	accessToken, refreshToken, err := auth.RefreshTokens(resfreshTokenRequest.Token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	existingToken.Token = refreshToken

	if err := database.DB.Model(&existingToken).Updates(existingToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to update refresh token in the database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}
