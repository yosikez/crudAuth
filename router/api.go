package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yosikez/crudAuth/controller"
	"github.com/yosikez/crudAuth/middleware"
)

func RegisterRoute(router *gin.Engine) {
	authController := controller.NewAuthController()

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.POST("/refresh-token", authController.RefreshToken)


	protected := router.Group("/api", middleware.AuthMiddleware())

	protected.GET("/hello",func(c *gin.Context) {
		username, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{
			"message" : "hello",
			"name" : username,
		})
	})
}