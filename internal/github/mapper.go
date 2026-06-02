package github

import "github.com/NemanjaTomic57/commitflow/internal/kafka"

func (c commitResponse) ToGitCommit() kafka.GitCommit {
	return kafka.GitCommit{
		ID:          c.SHA,
		AuthorName:  c.Commit.Author.Name,
		AuthorEmail: c.Commit.Author.Email,
		Message:     c.Commit.Message,
		CreatedAt:   c.Commit.Author.Date,
		URL:         c.HTMLURL,
		Provider:    "github",
	}
}
