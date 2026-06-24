CREATE TABLE IF NOT EXISTS git_commits (
  provider TEXT NOT NULL,
  id TEXT NOT NULL,
  path TEXT NOT NULL,
  path_with_namespace TEXT NOT NULL,
  author_name TEXT NOT NULL,
  author_email TEXT NOT NULL,
  message TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  ingested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  url TEXT NOT NULL,

  PRIMARY KEY (provider, id)
);

CREATE INDEX IF NOT EXISTS idx_git_commits_created_at
ON git_commits(created_at);

CREATE INDEX IF NOT EXISTS idx_git_commits_author_email
ON git_commits(author_email);
