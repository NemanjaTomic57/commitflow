package kafka

import (
	"encoding/json"
	"log"
	"os"

	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var topic = "git.commits"
var messages = make(chan []byte)

func Bootstrap() {
	go gitlab.GetAllCommits(messages)

	producer := newProducer()
	defer producer.Close()

	for message := range messages {
		produceKafkaEvents[gitlab.GitlabCommit](producer, message, topic)
	}

	producer.Flush(15 * 1000)
}

func newProducer() *kafka.Producer {
	bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		log.Fatalln("kafka.newProducer() -> KAFKA_BOOTSTRAP_SERVER is not set")
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServer})
	if err != nil {
		log.Fatalln("kafka.newProducer() -> error creating Kafka producer: %w", err)
	}

	return p
}

func produceKafkaEvents[T gitlab.GitAPIResponse](p *kafka.Producer, resp []byte, topic string) {
	var object []T

	err := json.Unmarshal(resp, &object)
	if err != nil {
		log.Println("produceKafkaEvents() -> error unmarshalling JSON: %w", err)
	}

	// Get results back from producing to Kafka and print to console
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Println("produceKafkaEvents() -> delivery failed:", ev.TopicPartition)
				} else {
					log.Println("produceKafkaEvents() -> delivered message to:", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce to Kafka topic
	for _, project := range object {
		projectBytes, err := json.Marshal(project)
		if err != nil {
			log.Printf("produceKafkaEvents() -> error marshalling project: %v", err)
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
