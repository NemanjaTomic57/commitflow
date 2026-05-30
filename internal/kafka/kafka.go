package kafka

import (
	"encoding/json"
	"log"
	"os"

	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

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

func ProduceKafkaEvents[T gitlab.GitAPIResponse](p *kafka.Producer, resp []byte, topic string) {
	var object []T

	err := json.Unmarshal(resp, &object)
	if err != nil {
		log.Println("ProduceKafkaEvents() -> error unmarshalling JSON: %w", err)
	}

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
	for _, project := range object {
		projectBytes, err := json.Marshal(project)
		if err != nil {
			log.Printf("ProduceKafkaEvents() -> error marshalling project: %v", err)
			continue
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          projectBytes,
		}, nil)
		if err != nil {
			log.Println("ProduceKafkaEvents() -> error producing message:", err)
		}
	}
}
