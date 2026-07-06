package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/NemanjaTomic57/commitflow/internal/github"
	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/NemanjaTomic57/commitflow/proto"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

const topic = "git_commits"

// Fetch all historical commits from Git.
func bootstrap(producer *ckafka.Producer) {
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

	for msg := range messages {
		kafka.ProduceEvent(producer, msg, topic)
	}
}

// Fetch commits from the last ten minutes only.
func pollGitAPI(producer *ckafka.Producer) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		log.Println("LOG pollGitAPI() -> starting to poll")

		var wg sync.WaitGroup
		messages := make(chan *proto.GitCommit)

		wg.Go(func() {
			gitlab.GetLastCommits(messages)
		})

		wg.Go(func() {
			github.GetLastCommits(messages)
		})

		go func() {
			wg.Wait()
			close(messages)
		}()

		for msg := range messages {
			kafka.ProduceEvent(producer, msg, topic)
		}

		<-ticker.C
	}
}

func handleDeliveryReports(producer *ckafka.Producer) {
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *ckafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Println("ERROR handleDeliveryReports() -> delivery failed:", ev.TopicPartition)
			} else {
				log.Println("LOG handleDeliveryReports() -> delivered message to:", ev.TopicPartition)
			}
		}
	}
}

func main() {
	_ = godotenv.Load()

	bootstrapFlag := flag.Bool("bootstrap", false, "Bootstrap historical Git data")
	flag.Parse()

	producer := kafka.NewProducer()
	defer func() {
		producer.Flush(1000 * 30)
		producer.Close()
	}()

	// Get results back from producing to Kafka and print to console
	go handleDeliveryReports(producer)

	if *bootstrapFlag {
		bootstrap(producer)
		log.Println("LOG bootstrapping historical data from Git is finished")
	}

	pollGitAPI(producer)
}
