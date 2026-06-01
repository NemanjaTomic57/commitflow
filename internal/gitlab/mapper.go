package gitlab

import "github.com/NemanjaTomic57/commitflow/internal/kafka"

func (commit gitlabCommit) ToGitCommit() kafka.GitCommit {
	return kafka.GitCommit{
		ID:          commit.ID,
		AuthorName:  commit.AuthorName,
		AuthorEmail: commit.AuthorEmail,
		Message:     commit.Message,
		CreatedAt:   commit.CreatedAt,
		URL:         commit.WebURL,
		Provider:    "gitlab",
	}
}
