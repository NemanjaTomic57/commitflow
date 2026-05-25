package main

import (
	"fmt"
	"log"

	"github.com/NemanjaTomic57/commitflow/internal/gitlab"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/NemanjaTomic57/commitflow/internal/utils"
	"github.com/joho/godotenv"
)

// TODO: Error handling on all functions
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("main() -> error loading .env file")
	}

	producer := kafka.NewProducer()
	defer producer.Close()

	projectIDs := gitlab.FetchProjectIDs()

	// TODO: Extract most of this to package gitlab and use concurrency
	for _, id := range projectIDs {
		url := baseURL + fmt.Sprintf("/projects/%d/repository/commits", id)
		topic := "git.commits"

		for url != "" {
			resp := gitlab.FetchAPI(url)
			url = gitlab.GetNextLink(resp)
			body := utils.ExtractBodyFromResponse(resp)
			resp.Body.Close()
			kafka.ProduceKafkaEvents[gitlab.GitlabCommit](producer, body, topic)
		}

	}
	producer.Flush(15 * 1000)
}
