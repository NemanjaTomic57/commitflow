package kafka

import (
	"log"

	"github.com/NemanjaTomic57/commitflow/proto"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	gproto "google.golang.org/protobuf/proto"
)

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
