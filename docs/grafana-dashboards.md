# Grafana Dashboards

Currently, we have set up dashboards for the following queries.

## GitCommits

### Pie Chart: Count per Provider

Title: CommitsPerProvider
Desc: Number of commits per git provider

```sql
SELECT provider, COUNT(*) FROM git_commits GROUP BY provider ORDER BY count DESC;
```

### Pie Chart: Count per Repository

Title: CommitsPerRepository
Desc: Number of commits per git repository

```sql
SELECT path, provider, count(*) FROM git_commits GROUP BY path, provider ORDER BY count DESC;
```

### Time Series: Commits over Time

Title: CommitsOverTime
Desc: Number of commits over time

```sql
SELECT created_at, SUM(COUNT(*)) OVER (ORDER BY created_at) FROM git_commits GROUP BY created_at;

SELECT created_at, SUM(COUNT(*)) OVER (ORDER BY created_at) FROM git_commits 
    WHERE path_with_namespace = 'NemanjaTomic57/commitflow' GROUP BY created_at;
```
