package gitlab

import (
	"encoding/json"
	"fmt"
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

func GetAllCommits(resp chan []byte) {
	// Fetch all project IDs
	projectIDs := fetchProjectIDs()

	// Interate through project IDs and fetch commits for each project
	for _, id := range projectIDs {
		url := fmt.Sprintf("%s/projects/%d/repository/commits", baseURL, id)
		fetchAPI(url, resp)
	}

	// Close channel when done
	close(resp)
}

func fetchAPI(url string, resp chan []byte) {
	// Iterate as long as there is an URL
	for url != "" {
		httpResponse := makeRequest(url)
		// Get the next link from the paginated result
		url = getNextLink(httpResponse)
		resp <- utils.ExtractBodyFromHTTPResponse(httpResponse)
		httpResponse.Body.Close()
	}
}

func makeRequest(url string) *http.Response {
	gitlabPAT := os.Getenv("GITLAB_PAT")
	if gitlabPAT == "" {
		log.Fatalln("GITLAB_PAT is not set")
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("gitlab.makeRequest() -> error creating the request: %w", err)
	}

	req.Header.Add("PRIVATE-TOKEN", gitlabPAT)

	// Send request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("gitlab.makeRequest() -> error sending the request:", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("gitlab.makeRequest() -> request status code error: %s with URL: %s", resp.Status, url)
	}

	return resp
}

// Extracts the next link from paginated http response
func getNextLink(resp *http.Response) string {
	linkHeader := resp.Header.Get("Link")

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

func fetchProjectIDs() []int {
	var resp = make(chan []byte)
	url := baseURL + "/projects?owned=true&per_page=1"

	// Fetch all projects and close channel when done
	go func() {
		fetchAPI(url, resp)
		close(resp)
	}()

	var projects []GitlabProject
	var projectIDs []int

	for r := range resp {
		// Unmarshal the paginated responses into objects
		json.Unmarshal(r, &projects)

		// Extract the IDs for each project in the resoponse
		for _, project := range projects {
			projectIDs = append(projectIDs, project.ID)
		}
	}

	// Return the project IDs
	return projectIDs
}
