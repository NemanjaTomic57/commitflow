package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumeMessage(consumer *kafka.Consumer, topic string) {
	err := consumer.Subscribe(topic, nil)
	if err != nil {
		log.Printf("kafka.ConsumeMessage() -> ERROR when subscribing to topic: %v\n", err)
	}

	for {
		event := consumer.Poll(100)
		switch e := event.(type) {
		case *kafka.Message:
			fmt.Println(e)
		}
	}
}
