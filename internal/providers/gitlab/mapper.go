package gitlab

import "github.com/NemanjaTomic57/commitflow/internal/kafka"

func (c commit) ToGitCommit() kafka.GitCommit {
	return kafka.GitCommit{
		ID:          c.ID,
		AuthorName:  c.AuthorName,
		AuthorEmail: c.AuthorEmail,
		Message:     c.Message,
		CreatedAt:   c.CreatedAt,
		URL:         c.WebURL,
		Provider:    "gitlab",
	}
}
