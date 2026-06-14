package main

import (
	"flag"
	"log"
	"sync"

	"github.com/NemanjaTomic57/commitflow/internal/aws"
	"github.com/NemanjaTomic57/commitflow/internal/github"
	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/joho/godotenv"
)

func bootstrap() {
	messages := make(chan kafka.GitCommit)
	var wg sync.WaitGroup

	wg.Go(func() {
		gitlab.GetAllCommits(messages)
	})

	wg.Go(func() {
		github.GetAllCommits(messages)
	})

	go func() {
		wg.Wait()
		close(messages)
	}()

	producer := kafka.NewProducer()
	defer producer.Close()

	topic := "git.commits"

	for message := range messages {
		// kafka.ProduceKafkaEvents(producer, message, topic)
		kafka.ProduceSchema(producer, message, topic)
	}

	producer.Flush(15 * 1000)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("main() -> error loading .env file")
	}

	bootstrapFlag := flag.Bool("bootstrap", false, "Bootstrap infrastructure")
	flag.Parse()

	// Fetch data data from Git if bootstrapping
	if *bootstrapFlag {
		err := aws.ResetS3Data()
		if err != nil {
			log.Fatal(err)
		}

		bootstrap()
	}

	// TODO: Implement cronjobs for API requests to Git
}
