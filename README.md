# CommitFlow

**CommitFlow** transforms Git events into actionable engineering intelligence using Kafka-powered streaming and cloud-native analytics.

The project explores how modern event-driven systems are built end-to-end with Apache Kafka — from ingesting external API data to processing, storage, and visualization.

The central question driving this project is:

> **How can scalable, production-ready data streaming systems be implemented with Kafka from end to end?**

CommitFlow implements a realistic streaming pipeline that collects development activity from GitHub and GitLab, processes it through Kafka, and delivers insights through analytics and dashboards.

## Architecture Goals

The goal is not only to use Kafka, but to understand how distributed streaming systems are designed, deployed, and operated in production environments.

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

The development environment will use **Vagrant** to provide:

- Kafka cluster
- Kafka Connect
- Supporting services
- Local observability stack

This setup allows rapid experimentation and local testing.

### Production

The production environment will be deployed on **AWS** using Terraform and Ansible for infrastructure automation.

The long-term goal is to create a reproducible and scalable streaming platform suitable for real-world workloads.

## Tech Stack

- **Apache Kafka**
- **Kafka Connect**
- **librdkafka**
- **Vagrant**
- **Terraform**
- **Ansible**
- **Grafana**
- **AWS**

## Learning Objectives

This project aims to build practical experience with:

- Event-driven architecture
- Distributed streaming systems
- Infrastructure as Code (IaC)
- Cloud-native deployment strategies
- Scalable data pipelines
- Kafka operations and observability

## Project Status

🚧 Work in progress — the architecture and infrastructure are actively being designed and implemented.
