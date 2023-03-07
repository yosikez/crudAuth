package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yosikez/crudAuth/config"
	"github.com/yosikez/crudAuth/controller"
	"github.com/yosikez/crudAuth/middleware"
)

func RegisterRoute(router *gin.Engine, conn *config.RabbitMQConnection, rmqCfg *config.RabbitMQ) {
	
	authController := controller.NewAuthController()
	todoController := controller.NewTodoController(conn, rmqCfg)


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

	protected.GET("/todos", todoController.FindAll)
	protected.GET("/todos/:id", todoController.FindById)
	protected.POST("/todos", todoController.Create)
	protected.PUT("/todos/:id", todoController.Update)
	protected.DELETE("/todos/:id", todoController.Delete)
}