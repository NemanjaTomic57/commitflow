package gitlab

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NemanjaTomic57/commitflow/internal/utils"
)

type GitlabProjectNamespace struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
	FullPath string `json:"full_path"`
	WebURL   string `json:"web_url"`
}

type GitlabProject struct {
	ID                     int                    `json:"id"`
	Description            string                 `json:"description"`
	PathWithNamespace      string                 `json:"path_with_namespace"`
	CreatedAt              string                 `json:"created_at"`
	WebURL                 string                 `json:"web_url"`
	GitlabProjectNamespace GitlabProjectNamespace `json:"namespace"`
}

type GitlabCommit struct {
	ID               string              `json:"id"`
	ShortID          string              `json:"short_id"`
	CreatedAt        time.Time           `json:"created_at"`
	ParentIDs        []string            `json:"parent_ids"`
	Title            string              `json:"title"`
	Message          string              `json:"message"`
	AuthorName       string              `json:"author_name"`
	AuthorEmail      string              `json:"author_email"`
	AuthoredDate     time.Time           `json:"authored_date"`
	CommitterName    string              `json:"committer_name"`
	CommitterEmail   string              `json:"committer_email"`
	CommittedDate    time.Time           `json:"committed_date"`
	Trailers         map[string]string   `json:"trailers"`
	ExtendedTrailers map[string][]string `json:"extended_trailers"`
	WebURL           string              `json:"web_url"`
}

type GitAPIResponse interface {
	GitlabProject | GitlabCommit
}

var baseURL = "https://gitlab.com/api/v4"

func FetchAPI(url string) *http.Response {
	gitlabPAT := os.Getenv("GITLAB_PAT")
	if gitlabPAT == "" {
		log.Fatal("GITLAB_PAT is not set")
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("fetchGitlabAPI() -> error creating the request:", err)
	}

	req.Header.Add("PRIVATE-TOKEN", gitlabPAT)

	// Send request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("fetchGitlabAPI() -> error sending the request:", err)
	}

	return resp
}

func GetNextLink(resp *http.Response) string {
	// Extract the next link from pagination
	linkHeader := resp.Header.Get("Link")

	if linkHeader == "" {
		return ""
	}

	for link := range strings.SplitSeq(linkHeader, ",") {
		parts := strings.Split(link, ";")

		if len(parts) < 2 {
			continue
		}

		urlPart := strings.TrimSpace(parts[0])
		relPart := strings.TrimSpace(parts[1])

		if relPart == `rel="next"` {
			return strings.Trim(urlPart, "<>")
		}
	}

	return ""
}

func FetchProjectIDs() []int {
	var projects []GitlabProject
	var projectIDs []int
	url := baseURL + "/projects?owned=true"

	for url != "" {
		resp := FetchAPI(url)
		url = GetNextLink(resp)
		body := utils.ExtractBodyFromResponse(resp)
		resp.Body.Close()

		json.Unmarshal(body, &projects)
		for _, project := range projects {
			projectIDs = append(projectIDs, project.ID)
		}
	}

	return projectIDs
}
