package main

import (
	"flag"
	"log"
	"sync"

	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/NemanjaTomic57/commitflow/internal/s3"
	"github.com/NemanjaTomic57/commitflow/proto"
	"github.com/joho/godotenv"
)

var topic = "git_commits"

func bootstrap() {
	messages := make(chan *proto.GitCommit)
	var wg sync.WaitGroup

	err := s3.ResetS3Data()
	if err != nil {
		log.Fatal(err)
	}

	wg.Go(func() {
		gitlab.GetAllCommits(messages)
	})

	// wg.Go(func() {
	// 	github.GetAllCommits(messages)
	// })

	go func() {
		wg.Wait()
		close(messages)
	}()

	producer := kafka.NewProducer()
	defer producer.Close()

	for message := range messages {
		go kafka.ProduceKafkaEvents(producer, message, topic)
	}

	producer.Flush(15 * 1000)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("main() -> error loading .env file")
	}

	consumer := kafka.NewConsumer()
	go kafka.ConsumeMessage(consumer, topic)
	defer consumer.Close()

	bootstrapFlag := flag.Bool("bootstrap", false, "Bootstrap infrastructure")
	flag.Parse()

	// Fetch data data from Git if bootstrapping
	if *bootstrapFlag {
		bootstrap()
	}

	// var wg sync.WaitGroup
	// wg.Add(1)
	// wg.Wait()
	// TODO: Implement cronjobs for API requests to Git
}
