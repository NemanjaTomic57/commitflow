package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/NemanjaTomic57/commitflow/internal/aws"
	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("main() -> error loading .env file")
	}

	bootstrap := flag.Bool("bootstrap", false, "Bootstrap infrastructure")
	flag.Parse()

	// Fetch data data from Git if bootstrapping
	if *bootstrap {
		err := aws.ResetS3Data()
		if err != nil {
			log.Fatal(err)
		}

		kafka.Bootstrap()
	}

	fmt.Println("TODO: Implement cronjobs for API requests to Git")
}
