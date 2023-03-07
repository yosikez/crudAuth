package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yosikez/crudAuth/database"
	"github.com/yosikez/crudAuth/helper/validation"
	"github.com/yosikez/crudAuth/router"
	"github.com/yosikez/crudAuth/rabbitmq"

)

func main() {
	// database
	err := database.Connect()
	if err != nil {
		panic(err)
	}

	// rabbitmq
	rmqCfg, rmq, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq : %v", err)
	}

	defer rmq.Connection.Close()
	defer rmq.Channel.Close()

	err = rmq.Channel.ExchangeDeclare(
		rmqCfg.ExchangeName,
		rmqCfg.ExchangeKind,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("failed to declare exchange : %v", err)
	}

	// declare gin.Engine
	r := gin.Default()
	// register the route
	router.RegisterRoute(r, rmq, rmqCfg)
	// register the custom validation
	validation.RegisterCustomValidation()
	// run the server on port 8000
	err = r.Run(":8000")
	
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	fmt.Println("server started on port 8000")

}