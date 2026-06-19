package gitlab

import (
	"github.com/NemanjaTomic57/commitflow/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c commit) ToGitCommit() *proto.GitCommit {
	provider := "gitlab"

	return proto.GitCommit_builder{
		Id:          &c.ID,
		AuthorName:  &c.AuthorName,
		AuthorEmail: &c.AuthorEmail,
		Message:     &c.Message,
		CreatedAt:   timestamppb.New(c.CreatedAt),
		Url:         &c.WebURL,
		Provider:    &provider,
	}.Build()
}
