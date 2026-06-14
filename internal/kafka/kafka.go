package kafka

import (
	"encoding/json"
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

func NewProducer() *kafka.Producer {
	bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		log.Fatalln("kafka.NewProducer() -> KAFKA_BOOTSTRAP_SERVER is not set")
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServer})
	if err != nil {
		log.Fatalln("kafka.NewProducer() -> error creating Kafka producer: %w", err)
	}

	return p
}

func ProduceKafkaEvents[T GitAPIResponse](p *kafka.Producer, message T, topic string) {
	// Get results back from producing to Kafka and print to console
	go func() {
		for e := range p.Events() {
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

	// Produce to Kafka topic
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("ProduceKafkaEvents() -> error marshalling project: %v", err)
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageBytes,
	}, nil)
	if err != nil {
		log.Println("ProduceKafkaEvents() -> error producing message:", err)
	}
}

// HACK: Use Kafka Connect built in JSON schema for PoC
type GitCommitMessage struct {
	Schema  Schema    `json:"schema"`
	Payload GitCommit `json:"payload"`
}

var GitCommitSchema = Schema{
	Type:     "struct",
	Optional: false,
	Name:     "commitflow.gitcommit",
	Fields: []SchemaField{
		{Type: "string", Optional: false, Field: "id"},
		{Type: "string", Optional: false, Field: "author_name"},
		{Type: "string", Optional: false, Field: "author_email"},
		{Type: "string", Optional: false, Field: "message"},
		{Type: "string", Optional: false, Field: "created_at"},
		{Type: "string", Optional: false, Field: "url"},
		{Type: "string", Optional: false, Field: "provider"},
	},
}

type Message struct {
	Schema  Schema  `json:"schema"`
	Payload Payload `json:"payload"`
}

type Schema struct {
	Type     string        `json:"type"`
	Fields   []SchemaField `json:"fields"`
	Optional bool          `json:"optional"`
	Name     string        `json:"name"`
}

type SchemaField struct {
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Field    string `json:"field"`
}

type Payload struct {
	RegisterTime int64  `json:"registertime"`
	UserID       string `json:"userid"`
	RegionID     string `json:"regionid"`
	Gender       string `json:"gender"`
}

func ProduceSchema(p *kafka.Producer, message GitCommit, topic string) {
	// Get results back from producing to Kafka and print to console
	go func() {
		for e := range p.Events() {
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

	event := GitCommitMessage{
		Schema:  GitCommitSchema,
		Payload: message,
	}

	// Produce to Kafka topic
	messageBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("ProduceKafkaEvents() -> error marshalling project: %v", err)
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageBytes,
	}, nil)
	if err != nil {
		log.Println("ProduceKafkaEvents() -> error producing message:", err)
	}
}
