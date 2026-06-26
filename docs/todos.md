# To Do's

- Implement cronjobs for recurring API requests
- Scale the application to several topics
- Automate the Grafana dashboards

## Production

- Standalone EC2 instance with Docker compose
- Deployment to AWS with Terraform
- GitOps for updating images on Docker

## Grafana Dashboards

Currently, we have set up dashboards for the following queries.

### GitCommits

### Pie Chart: Count per Provider

SELECT provider, COUNT(*) FROM git_commits GROUP BY provider ORDER BY count DESC;

### Pie Chart: Count per Repository

SELECT path, provider, count(*) FROM git_commits GROUP BY path, provider ORDER BY count DESC;

### Time Series: Commits / Time

SELECT created_at, SUM(COUNT(*)) OVER (ORDER BY created_at) FROM git_commits GROUP BY created_at;
SELECT created_at, SUM(COUNT(*)) OVER (ORDER BY created_at) FROM git_commits WHERE path_with_namespace = 'NemanjaTomic57/commitflow' GROUP BY created_at;
