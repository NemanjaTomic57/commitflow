package main

import (
	"flag"
	"sync"

	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/NemanjaTomic57/commitflow/proto"
	"github.com/joho/godotenv"
)

var topic = "git_commits"

func bootstrap() {
	var wg sync.WaitGroup
	messages := make(chan *proto.GitCommit)

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
		go kafka.ProduceEvent(producer, message, topic)
	}

	producer.Flush(15 * 1000)
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
