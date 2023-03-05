package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yosikez/crudAuth/database"
	"github.com/yosikez/crudAuth/helper/validation"
	"github.com/yosikez/crudAuth/router"
)

func main() {
	err := database.Connect()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	router.RegisterRoute(r)
	validation.RegisterCustomValidation()

	err = r.Run(":8000")
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	fmt.Println("server started on port 8000")

}