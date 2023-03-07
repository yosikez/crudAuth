package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yosikez/crudAuth/config"
	"github.com/yosikez/crudAuth/database"
	"github.com/yosikez/crudAuth/model"
	amqp "github.com/rabbitmq/amqp091-go"
	cusMessage "github.com/yosikez/custom-error-message"
)

type TodoController struct {
	rmq *config.RabbitMQConnection
	rmqCfg *config.RabbitMQ
}

func NewTodoController(rqConnection *config.RabbitMQConnection, rqConfig *config.RabbitMQ) *TodoController {
	return &TodoController{
		rmq: rqConnection,
		rmqCfg: rqConfig,
	}
}

func (t *TodoController) FindAll(c *gin.Context) {
	var todos []model.Todo

	userId := c.GetUint("userId")

	result := database.DB.Where("user_id = ?", userId).Find(&todos)

	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find all todo",
			"error":   err.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "failed to find all todo",
			"error":   "record not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})
}

func (t *TodoController) FindById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	userId := c.GetUint("userId")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid todo id",
			"error":   "id must be a number",
		})
		return
	}

	var todo model.Todo
	if err := database.DB.Where("user_id = ?", userId).First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "failed to find todo",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": todo,
	})
}

func (t *TodoController) Create(c *gin.Context) {
	userId := c.GetUint("userId")
	var todo model.Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		errFields := cusMessage.GetErrMess(err, todo, nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid validation",
			"errors":  errFields,
		})
		return
	}

	todo.UserId = userId

	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create todo",
			"error":   err.Error(),
		})
		return
	}

	q, err := t.rmq.Channel.QueueDeclare("todo_create_queue", false, false, false, false, nil,)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to declare queue rabbitmq",
			"error":   err.Error(),
		})
		return
	}

	err = t.rmq.Channel.QueueBind("todo_create_queue", "todo_create_queue", t.rmqCfg.ExchangeName, false, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to bind a queue",
			"error" : err.Error(),
		})
		return
	}

	msg, err := json.Marshal(todo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to marshal json for message rabbitmq",
			"error":   err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = t.rmq.Channel.PublishWithContext(ctx, t.rmqCfg.ExchangeName, q.Name, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to publish message to rabbitmq",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": todo,
	})
}

func (t *TodoController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : "invalid todo id",
			"error": "id must be a number",
		})
		return
	}

	var todo model.Todo
	if err := database.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : "failed to find todo to update",
			"error" : err.Error(),
		})
		return
	}

	if err := c.ShouldBind(&todo); err != nil {
		errFields := cusMessage.GetErrMess(err, todo, nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : "validation error",
			"error" : errFields,
		})
		return
	}

	if err := database.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to create todo",
			"error": err.Error(),
		})
		return
	}

	q, err := t.rmq.Channel.QueueDeclare("todo_update_queue", false, false, false, false, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to declare queue rabbitmq",
			"error":   err.Error(),
		})
		return
	}

	err = t.rmq.Channel.QueueBind("todo_update_queue", "todo_update_queue", t.rmqCfg.ExchangeName, false, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to bind a queue",
			"error" : err.Error(),
		})
		return
	}

	msg, err := json.Marshal(&todo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to marshal json for msg rabbitmq",
			"error" : err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = t.rmq.Channel.PublishWithContext(ctx, t.rmqCfg.ExchangeName, q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: msg,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to publish message to rabbitmq",
			"error" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data" : todo,
	})
}

func (t *TodoController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : "invalid todo id",
			"error" : err.Error(),
		})
		return
	}

	var todo model.Todo
	if err := database.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message" : "failed to find todo to delete",
			"error" : err.Error(),
		})
		return
	}

	if err := database.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to delete todo",
			"error" : err.Error(),
		})
		return
	}

	q, err := t.rmq.Channel.QueueDeclare("todo_delete_queue", false, false, false, false, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to declare queue rabbitmq",
			"error" : err.Error(),
		})
		return
	}

	err = t.rmq.Channel.QueueBind("todo_delete_queue", "todo_delete_queue", t.rmqCfg.ExchangeName, false, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to bind a queue",
			"error" : err.Error(),
		})
		return
	}

	msg, err := json.Marshal(&todo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to marshal json for the message rabbitmq",
			"error" : err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = t.rmq.Channel.PublishWithContext(ctx, t.rmqCfg.ExchangeName, q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: msg,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "failed to publish message to rabbitmq",
			"error" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "todo deleted successfully",
	})

}