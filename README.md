# CommitFlow

The project is a pilot project focused on providing a PoC for Kafka and Kafka Connect.

## What This Project Is (... And What It Is Not)

The goal of the project is to get a thorough understanding of the Kafka ecosystem, specifically how it works together with external data sources and destinations. The main focus is to grasp the underlying technology in its fullest by implementing a real-world use case. The leading question is the following:

**How are scalable solutions for data streaming implemented in Kafka end to end?**

Principles like high availability are not part of the requirements, as speed and flexibility are of greater importance.

## Requirements

We will create a Kafka cluster in a virtual machine with the following requirements:

- Configure a scalable Kafka infrastructure with Ansible and Terraform
- Fetch data from GitHub and GitLab APIs with cron jobs
- Push the data to Kafka topics with the librdkafka library
- Use Kafka Connect to sink the data to S3
- Create a dashboard in Grafana to display metrics

We plan to create a development environment using Docker Compose. Later, we'll set up a production environment on Hetzner Cloud with Docker for the runtime environment. 
