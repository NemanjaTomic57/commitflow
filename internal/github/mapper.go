package github

import (
	"github.com/NemanjaTomic57/commitflow/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c commitResponse) ToGitCommit(repo repository) *proto.GitCommit {
	provider := "github"

	return proto.GitCommit_builder{
		Provider:          &provider,
		Id:                &c.SHA,
		Path:              &repo.Name,
		PathWithNamespace: &repo.FullName,
		AuthorName:        &c.Commit.Author.Name,
		AuthorEmail:       &c.Commit.Author.Email,
		Message:           &c.Commit.Message,
		Url:               &c.HTMLURL,
		CreatedAt:         timestamppb.New(c.Commit.Author.Date),
	}.Build()
}
