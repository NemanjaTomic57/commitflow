# CommitFlow

**CommitFlow** implements a realistic streaming pipeline that collects development activity from GitHub and GitLab, processes it through Kafka, and delivers insights through analytics and dashboards.

The project explores how modern event-driven systems are built end-to-end with Apache Kafka — from ingesting external API data to processing, storage, and visualization.

The central question driving this project is:

> **How can scalable, production-ready data streaming systems be implemented with Kafka from end to end?**

The goal is not only to use Kafka, but to understand how distributed streaming systems are designed, deployed, and operated in production environments.

## Getting Started

### Prerequisites

* Docker Compose
* Go
* GitHub Personal Access Token with permissions `repo`
* GitLab Personal Access Token with permissions `api`

Create your personal access tokens:

* GitHub: https://github.com/settings/tokens
* GitLab: https://gitlab.com/-/user_settings/personal_access_tokens

Copy the example environment file and add both tokens:

```bash
cp .env.example .env
```

### Start the Development Environment

Start Kafka, PostgreSQL, and Grafana:

```bash
docker compose up -d
```

### Bootstrap Historical Data

Start the consumer first to initialize the database and begin consuming Kafka topics:

```bash
go run ./cmd/consumer/consumer.go
```

In a second terminal, run the producer once with the bootstrap flag to import your complete GitHub and GitLab history:

```bash
go run ./cmd/producer/producer.go --bootstrap
```

Within about a minute, your data should appear in Grafana.

### Access Services

| Service    | URL / Command                                 | Credentials             |
| ---------- | --------------------------------------------- | ----------------------- |
| Grafana    | http://localhost:3000                         | `admin` / `password`    |
| PostgreSQL | `psql -U postgres -h localhost -d commitflow` | `postgres` / `password` |

## Architecture Goals

### Infrastructure & Deployment

- Provision infrastructure using **Terraform**
- Configure services with **Ansible**
- Deploy scalable Kafka components in a cloud environment
- Support both local development and production deployments

### Data Ingestion

- Fetch repository activity from:
  - GitHub API
  - GitLab API
- Schedule ingestion jobs using cron-based workers
- Publish events to Kafka topics using **librdkafka**

### Streaming & Storage

- Design topic structures for scalable event processing
- Use **Kafka Connect** for data integration
- Sink streamed data into **AWS S3** and **PostgreSQL**

### Monitoring & Analytics

- Collect operational metrics with the collected data from Git
- Visualize pipeline health and engineering insights with **Grafana**
- Monitor Kafka infrastructure and throughput

## Environments

### Local Development

The development environment will use **Docker** to provide:

- Kafka 
- PostgreSQL
- Grafana Dashboards

This setup allows rapid experimentation and local testing.

### Production

The production environment will be deployed on **AWS** using Terraform and Ansible for infrastructure automation.

The long-term goal is to create a reproducible and scalable streaming platform suitable for real-world workloads.

## Learning Objectives

This project aims to build practical experience with:

- Event-driven architecture
- Distributed streaming systems
- Infrastructure as Code (IaC)
- Cloud-native deployment strategies
- Scalable data pipelines
- Kafka operations and observability
