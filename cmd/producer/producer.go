package main

import (
	"flag"
	"log"
	"sync"

	"github.com/NemanjaTomic57/commitflow/internal/github"
	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/NemanjaTomic57/commitflow/proto"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

var topic = "git_commits"

func bootstrap() {
	var wg sync.WaitGroup
	messages := make(chan *proto.GitCommit)

	wg.Go(func() {
		gitlab.GetAllCommits(messages)
	})

	wg.Go(func() {
		github.GetAllCommits(messages)
	})

	// Wait until all commits are fetched
	go func() {
		wg.Wait()
		close(messages)
	}()

	producer := kafka.NewProducer()

	// Get results back from producing to Kafka and print to console
	go handleDeliveryReports(producer)

	// Produce Kafka events for each message
	var produceWg sync.WaitGroup

	for message := range messages {
		produceWg.Go(func() {
			kafka.ProduceEvent(producer, message, topic)
		})
	}

	// Wait until producer has processed all commits
	produceWg.Wait()

	producer.Flush(1000 * 30)
	producer.Close()
}

func handleDeliveryReports(producer *ckafka.Producer) {
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *ckafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Println("ERROR kafka.ProduceKafkaEvents() -> delivery failed:", ev.TopicPartition)
			} else {
				log.Println("LOG kafka.ProduceKafkaEvents() -> delivered message to:", ev.TopicPartition)
			}
		}
	}
}

func main() {
	godotenv.Load()
	bootstrapFlag := flag.Bool("bootstrap", false, "Bootstrap infrastructure")
	flag.Parse()

	// Fetch data data from Git if bootstrapping
	if *bootstrapFlag {
		bootstrap()
	}
	// TODO: Implement cronjobs for API requests to Git
}
