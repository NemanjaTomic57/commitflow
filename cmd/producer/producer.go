package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/NemanjaTomic57/commitflow/internal/github"
	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/NemanjaTomic57/commitflow/proto"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

var topic = "git_commits"

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
	messages := make(chan *proto.GitCommit)
	defer close(messages)

	for {
		time.Sleep(10 * time.Second)

		log.Println("starting poll")

		start := time.Now()

		log.Println("calling GetLastCommits")
		go github.GetLastCommits(messages)
		log.Println("GetLastCommits returned after", time.Since(start))

		if len(messages) == 0 {
			log.Println("no messages received from git")
			continue
		}

		log.Println("reading channel")
		for msg := range messages {
			kafka.ProduceEvent(producer, msg, topic)
		}

		log.Println("channel closed")
	}
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

	producer := kafka.NewProducer()

	// Get results back from producing to Kafka and print to console
	go handleDeliveryReports(producer)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// bootstrap(producer)

	log.Println("LOG bootstrapping historical data from Git is finished")

	go pollGitAPI(producer)

	<-ctx.Done()

	producer.Flush(1000 * 30)
	producer.Close()
}
