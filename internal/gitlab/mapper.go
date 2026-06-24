package gitlab

import (
	"github.com/NemanjaTomic57/commitflow/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c commit) ToGitCommit(project project) *proto.GitCommit {
	provider := "gitlab"

	return proto.GitCommit_builder{
		Provider:          &provider,
		Id:                &c.ID,
		Path:              &project.Path,
		PathWithNamespace: &project.PathWithNamespace,
		AuthorName:        &c.AuthorName,
		AuthorEmail:       &c.AuthorEmail,
		Message:           &c.Message,
		Url:               &c.WebURL,
		CreatedAt:         timestamppb.New(c.CreatedAt),
	}.Build()
}
