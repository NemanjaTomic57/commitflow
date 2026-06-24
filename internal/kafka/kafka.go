package kafka

import (
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type GitCommit struct {
	ID          string    `json:"id"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	URL         string    `json:"url"`
	Provider    string    `json:"provider"`
}

type GitAPIResponse interface {
	GitCommit
}

func NewConsumer() *kafka.Consumer {
	bootstrapServer := getBootstrapServer()
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          "s3-sink",
	})
	if err != nil {
		log.Printf("kafka.NewConsumer() -> ERROR when creating consumer: %v\n", err)
	}

	return consumer
}

func getBootstrapServer() string {
	bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		log.Fatalln("kafka.NewProducer() -> KAFKA_BOOTSTRAP_SERVER is not set")
	}

	return bootstrapServer
}
