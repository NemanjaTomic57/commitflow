# CommitFlow

**CommitFlow** transforms Git events into actionable engineering intelligence using Kafka-powered streaming and cloud-native analytics.

The project explores how modern event-driven systems are built end-to-end with Apache Kafka — from ingesting external API data to processing, storage, and visualization.

---

## Overview

This project is designed as a hands-on deep dive into the Kafka ecosystem and surrounding infrastructure.

Instead of building a toy example, CommitFlow implements a realistic streaming pipeline that collects development activity from platforms like GitHub and GitLab, processes it through Kafka, and delivers insights through analytics and dashboards.

The central question driving this project is:

> **How can scalable, production-ready data streaming systems be implemented with Kafka from end to end?**

---

## Architecture Goals

The project focuses on understanding how Kafka integrates with:

- External APIs and event sources
- Producers and consumers
- Kafka Connect
- Object storage systems
- Monitoring and observability tooling
- Infrastructure automation and deployment workflows

The goal is not only to use Kafka, but to understand how distributed streaming systems are designed, deployed, and operated in production environments.

---

## Planned Features

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
- Sink streamed data into **Amazon S3-compatible storage**

### Monitoring & Analytics

- Collect operational metrics
- Visualize pipeline health and engineering insights with **Grafana**
- Monitor Kafka infrastructure and throughput

---

## Environments

### Local Development

The development environment will use **Docker Compose** to provide:

- Kafka cluster
- Kafka Connect
- Supporting services
- Local observability stack

This setup allows rapid experimentation and local testing.

### Production

The production environment will be deployed on **Hetzner Cloud** using Docker-based workloads and infrastructure automation.

The long-term goal is to create a reproducible and scalable streaming platform suitable for real-world workloads.

---

## Tech Stack

- **Apache Kafka**
- **Kafka Connect**
- **librdkafka**
- **Docker & Docker Compose**
- **Terraform**
- **Ansible**
- **Grafana**
- **Amazon S3-compatible storage**
- **Hetzner Cloud**

---

## Learning Objectives

This project aims to build practical experience with:

- Event-driven architecture
- Distributed streaming systems
- Infrastructure as Code (IaC)
- Cloud-native deployment strategies
- Scalable data pipelines
- Kafka operations and observability

---

## Project Status

🚧 Work in progress — the architecture and infrastructure are actively being designed and implemented.
