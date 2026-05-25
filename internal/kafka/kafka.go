package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer() *kafka.Producer {
	bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		log.Fatal("KAFKA_BOOTSTRAP_SERVER is not set")
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServer})
	if err != nil {
		log.Fatal("main() -> error creating Kafka producer: ", err)
	}

	return p
}

func ProduceKafkaEvents[T gitlab.GitAPIResponse](p *kafka.Producer, resp []byte, topic string) {
	var object []T

	err := json.Unmarshal(resp, &object)
	if err != nil {
		log.Fatal("produceKafkaEvents() -> error unmarshalling JSON:", err)
	}

	// Get results back from producing to Kafka and print to console
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("produceKafkaEvents() -> delivery failed %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("produceKafkaEvents() -> delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce to Kafka topic
	for _, project := range object {
		projectBytes, err := json.Marshal(project)
		if err != nil {
			log.Println("produceKafkaEvents() -> error marshalling project:", err)
			continue
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          projectBytes,
		}, nil)
		if err != nil {
			log.Println("produceKafkaEvents() -> error producing message:", err)
		}
	}
}
