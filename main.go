package main

import (
	"log"

	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/joho/godotenv"
)

var topic = "git.commits"
var messages = make(chan []byte)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("main() -> error loading .env file")
	}

	producer := kafka.NewProducer()
	if err != nil {
		log.Printf("main() -> error creating kafka producer: %v", err)
	}
	defer producer.Close()

	go gitlab.GetAllCommits(messages)

	for message := range messages {
		kafka.ProduceKafkaEvents[gitlab.GitlabCommit](producer, message, topic)
	}

	producer.Flush(15 * 1000)
}
