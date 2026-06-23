package kafka

import (
	"log"
	"os"
	"time"

	"github.com/NemanjaTomic57/commitflow/proto"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	gproto "google.golang.org/protobuf/proto"
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

func NewProducer() *kafka.Producer {
	bootstrapServer := getBootstrapServer()
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
	})
	if err != nil {
		log.Fatalln("kafka.NewProducer() -> error creating Kafka producer: %w", err)
	}

	return producer
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

func ProduceKafkaEvents(producer *kafka.Producer, message *proto.GitCommit, topic string) {
	// Get results back from producing to Kafka and print to console
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Println("ProduceKafkaEvents() -> delivery failed:", ev.TopicPartition)
				} else {
					log.Println("ProduceKafkaEvents() -> delivered message to:", ev.TopicPartition)
				}
			}
		}
	}()

	// fmt.Println(message)

	messageBytes, err := gproto.Marshal(message)
	if err != nil {
		log.Printf("ProduceKafkaEvents() -> error marshalling object: %v", err)
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: messageBytes,
	}, nil)
	if err != nil {
		log.Println("ProduceKafkaEvents() -> error producing message:", err)
	}
}

func getBootstrapServer() string {
	bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		log.Fatalln("kafka.NewProducer() -> KAFKA_BOOTSTRAP_SERVER is not set")
	}

	return bootstrapServer
}
