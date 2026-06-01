package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/NemanjaTomic57/commitflow/internal/aws"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/NemanjaTomic57/commitflow/internal/providers/gitlab"
	"github.com/joho/godotenv"
)

// Creates the GitHub base URL with the variable username
func baseURL() string {
	owner := os.Getenv("GITHUB_USERNAME")
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s", owner)
	return baseURL
}

func bootstrap() {
	messages := make(chan kafka.GitCommit)

	go gitlab.GetAllCommits(messages)

	producer := kafka.NewProducer()
	defer producer.Close()

	topic := "git.commits"

	for message := range messages {
		kafka.ProduceKafkaEvents(producer, message, topic)
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
