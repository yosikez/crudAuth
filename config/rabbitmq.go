package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Host         string
	Username     string
	Password     string
	ExchangeName string
	ExchangeKind string
}

type RabbitMQConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (rmq *RabbitMQ) GetUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s/", rmq.Username, rmq.Password, rmq.Host)
}

func LoadRabbitMQ() (*RabbitMQ, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("failed to load .env file")
		return nil, err
	}

	rabbitmqConfig := &RabbitMQ{
		Host:         os.Getenv("RMQ_HOST"),
		Username:     os.Getenv("RMQ_USERNAME"),
		Password:     os.Getenv("RMQ_PASSWORD"),
		ExchangeName: os.Getenv("RMQ_EXCHANGE_NAME"),
		ExchangeKind: os.Getenv("RMQ_EXCHANGE_KIND"),
	}

	return rabbitmqConfig, nil
}
