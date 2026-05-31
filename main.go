package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/NemanjaTomic57/commitflow/internal/aws"
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

		// kafka.Bootstrap()
		os.Exit(0)
	}

	fmt.Println("Nothing to do.")
}
