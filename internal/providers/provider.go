package providers

import (
	"net/http"

	"github.com/NemanjaTomic57/commitflow/internal/kafka"
)

type Provider interface {
	GetAllCommits(chan kafka.GitCommit)
	fetchProjectIDs() []int
	executeRequest(url string) *http.Response
	getNextLink(resp *http.Response) string
}
