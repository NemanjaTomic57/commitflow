package kafka

import (
	"log"
	"os"
	"strings"
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

var Topic = "git_commits"

func NewProducer() *kafka.Producer {
	bootstrapServer := getBootstrapServer()
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
	})
	if err != nil {
		log.Fatalf("ERROR kafka.NewProducer() -> creating Kafka producer failed: %w", err)
	}

	return producer
}

func ProduceEvent(producer *kafka.Producer, message *proto.GitCommit, topic string) {
	// Remove whitespace from commit messages
	message.SetMessage(strings.TrimSpace(message.GetMessage()))

	messageBytes, err := gproto.Marshal(message)
	if err != nil {
		log.Printf("ERROR kafka.ProduceKafkaEvents() -> marshalling object failed: %v", err)
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: messageBytes,
	}, nil)
	if err != nil {
		log.Println("ERROR kafka.ProduceKafkaEvents() -> producing message failed:", err)
	}
}

func NewConsumer(groupID string) *kafka.Consumer {
	bootstrapServer := getBootstrapServer()
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          groupID,
	})
	if err != nil {
		log.Printf("ERROR kafka.NewConsumer() -> create consumer failed: %v\n", err)
	}

	return consumer
}

func ConsumeEvent(consumer *kafka.Consumer, topic string, messages chan *proto.GitCommit) {
	err := consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("kafka.ConsumeEvent() -> ERROR when subscribing to topic: %v\n", err)
	}

	for {
		var message proto.GitCommit
		event, err := consumer.ReadMessage(time.Millisecond * 100)
		if err != nil {
			continue
		}

		if err := gproto.Unmarshal(event.Value, &message); err != nil {
			log.Printf("ERROR kafka.ConsumeEvent() -> failed to unmarshal message: %v", err)
			continue
		}

		messages <- &message
	}
}

func getBootstrapServer() string {
	bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		log.Fatalln("ERROR kafka.NewProducer() -> KAFKA_BOOTSTRAP_SERVER is not set")
	}

	return bootstrapServer
}
