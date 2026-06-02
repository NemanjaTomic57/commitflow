package gitlab

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NemanjaTomic57/commitflow/internal/kafka"
	"github.com/NemanjaTomic57/commitflow/internal/utils"
)

var baseURL = "https://gitlab.com/api/v4"

// Get all commits for all projects.
func GetAllCommits(messages chan kafka.GitCommit) {
	defer close(messages)

	// Fetch all project IDs
	projects := fetchAllProjects()

	// Interate through project IDs and fetch commits for each project
	for _, project := range projects {
		url := fmt.Sprintf("%s/projects/%d/repository/commits", baseURL, project.ID)
		commits := fetchAPI[commit](url)

		for _, commit := range commits {
			message := commit.ToGitCommit()
			messages <- message
		}
	}
}

// Fetches all projects for the authenticated user.
func fetchAllProjects() []project {
	url := baseURL + "/projects?membership=true&per_page=1"
	return fetchAPI[project](url)
}

// Fetches the API endpoint and returns JSON array. If the response is
// paginated, fetch all pages.
func fetchAPI[T responseType](url string) []T {
	var result []T

	// Iterate as long as there is an URL
	for url != "" {
		// Make the API request with the current page
		httpResponse := executeRequest(url)

		// Get the next link from the paginated result
		url = getNextLink(httpResponse)

		// Send []byte of http response to the channel
		pageResp := utils.ExtractBodyFromHTTPResponse(httpResponse)
		httpResponse.Body.Close()

		// Unmarshal into corresponding object array...
		var page []T
		if err := json.Unmarshal(pageResp, &page); err != nil {
			log.Printf("gitlab.fetchAPI() -> error unmarshalling response: %v", err)
			continue
		}
		// ...and append.
		result = append(result, page...)
	}

	return result
}

// Makes a single request to the GitLab API.
func executeRequest(url string) *http.Response {
	gitlabPAT := os.Getenv("GITLAB_PAT")
	if gitlabPAT == "" {
		log.Fatalln("GITLAB_PAT is not set")
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("gitlab.executeRequest() -> error creating the request: %v", err)
	}

	req.Header.Add("PRIVATE-TOKEN", gitlabPAT)

	// Send request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("gitlab.executeRequest() -> error sending the request: %v", err)
		return nil
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("gitlab.executeRequest() -> request status code error: %s with URL: %s", resp.Status, url)
		return nil
	}

	return resp
}

// Extracts the next link from the paginated GitLab API response.
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
