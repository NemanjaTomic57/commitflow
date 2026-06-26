package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/NemanjaTomic57/commitflow/proto"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

var query = `
INSERT INTO git_commits (
	provider,
	id,
	path,
	path_with_namespace,
	author_name,
	author_email,
	message,
	url,
	created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (provider, id) DO UPDATE SET
	path = EXCLUDED.path,
	path_with_namespace = EXCLUDED.path_with_namespace,
	author_name = EXCLUDED.author_name,
	author_email = EXCLUDED.author_email,
	message = EXCLUDED.message,
	url = EXCLUDED.url,
	created_at = EXCLUDED.created_at
`

func migrateDB() {
	connectionString := os.Getenv("POSTGRES_URL")

	m, err := migrate.New(
		"file://migrations",
		connectionString,
	)
	if err != nil {
		log.Fatalf("ERROR failed to create migrate instance: %v", err)
	}
	if err = m.Up(); err != migrate.ErrNoChange {
		if err != nil {
			log.Fatalf("ERROR applying database migration failed: %v", err)
		}
		log.Println("LOG database migration successful")
	} else {
		log.Println("LOG database already migrated")
	}
}

func postgresSink(ctx context.Context) {
	messages := make(chan *proto.GitCommit)

	consumer := kafka.NewConsumer("postgres-sink")
	defer consumer.Close()

	go kafka.ConsumeEvent(ctx, consumer, kafka.Topic, messages)

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("ERROR postgresSink() -> could not open database connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("ERROR postgresSink() -> could not establish database connection: %v", err)
	}

	for {
		message := <-messages
		log.Printf("LOG consuming git commit with repository path: %s, message: %v", message.GetPathWithNamespace(), message.GetMessage())

		_, err := db.ExecContext(
			ctx,
			query,
			message.GetProvider(),
			message.GetId(),
			message.GetPath(),
			message.GetPathWithNamespace(),
			message.GetAuthorName(),
			message.GetAuthorEmail(),
			message.GetMessage(),
			message.GetUrl(),
			message.GetCreatedAt().AsTime(),
		)
		if err != nil {
			log.Printf("ERROR postgresSink() -> failed to insert commit %s/%s: %v",
				message.GetProvider(),
				message.GetId(),
				err)
		}
	}
}

func main() {
	godotenv.Load()

	migrateDB()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go postgresSink(ctx)

	<-ctx.Done()
}
