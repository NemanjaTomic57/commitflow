package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go postgresSink(ctx)

	<-ctx.Done()
}

func postgresSink(ctx context.Context) {
	messages := make(chan string)

	consumer := kafka.NewConsumer("postgres-sink")
	defer consumer.Close()

	go kafka.ConsumeEvent(ctx, consumer, kafka.Topic, messages)

	for {
		fmt.Println(<-messages)
	}
}
