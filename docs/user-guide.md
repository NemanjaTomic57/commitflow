# User Guide

To run this application in your development environment, you will need to fulfill the following prerequisites:

You need to...

- ...install Git.
- ...install Docker Compose.
- ...install Golang.
- ...create a GitHub Personal Access Token.
- ...create a GitLab Personal Access Token.

## GitHub and GitLab Personal Access Token

We need a personal access token for GitLab and GitHub APIs. Create them in your account and add them to the project by copying the environment file example and filling out the placeholders.

### GitLab Personal Access Token

URL:

https://gitlab.com/-/user_settings/personal_access_tokens

Permissions to include:

- api

### GitHub Personal Access Token

URL:

https://github.com/settings/tokens

- repo

### Update .env file

Once you created and copied the PATs into your password manager, you can copy the .env.example into .env and add the PATs for Git into the .env file.

## Provision the environment with Docker

```bash
docker compose up -d
```

## PostgreSQL Database

User: postgres
Password: password

If you want to connect to the PostgreSQL database, you can do so from within the VM as the vagrant user.

```bash
psql -U postgres -h localhost -d commitflow
```

## Grafana Dashboard

The Grafana dashboard is available under http://localhost:3000.

User: admin
Password: password

## Bootstrap you data from Git

To bootstrap you data from your Git account, the first thing you have to do is to run the consumer. This will then connect to Postgres, migrate the database and consume from the Kafka topics until something produces to Kafka.

```bash
go run ./cmd/consumer/consumer.go
```

After you run the consumer, open a second terminal window and start the producer with the bootstrap flag. This will collect all the historical data from your GitHub and GitLab and produce corresponding messages to Kafka topics.

```bash
go run ./cmd/producer/producer.go --bootstrap
```

After you run both the consumer and producer, you should be able to see your data from Git in Grafana. It is worthwhile to take a look at the created dashboards under http://localhost:3000. The data should be populated within a minute or so.
