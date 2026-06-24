package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NemanjaTomic57/commitflow/internal/utils"
	"github.com/NemanjaTomic57/commitflow/proto"
)

var baseURL = "https://api.github.com"

// Get all commits for every project.
func GetAllCommits(messages chan *proto.GitCommit) {
	// Fetch all project IDs
	repositories := fetchAllProjects()

	// Interate through project IDs and fetch commits for each project
	for _, repo := range repositories {
		url := fmt.Sprintf("%s/repos/%s/%s/commits?per_page=100", baseURL, repo.Owner.Login, repo.Name)
		commits := fetchAPI[commitResponse](url)

		for _, commit := range commits {
			message := commit.ToGitCommit(repo)
			messages <- message
		}
	}
}

// Fetches all projects for the authenticated user.
func fetchAllProjects() []repository {
	url := baseURL + "/user/repos?per_page=1"
	return fetchAPI[repository](url)
}

// Fetches the API endpoint and returns JSON array. If the response is
// paginated, fetch all pages.
func fetchAPI[T responseType](url string) []T {
	var result []T

	// Iterate as long as there is an URL
	for url != "" {
		// Make the API request with the current page
		httpResponse, err := executeRequest(url)
		if err != nil {
			log.Println(err)
			continue
		}

		// Get the next link from the paginated result
		url = getNextLink(httpResponse)

		// Send []byte of http response to the channel
		pageResp := utils.ExtractBodyFromHTTPResponse(httpResponse)
		httpResponse.Body.Close()

		// Unmarshal into corresponding object array...
		var page []T
		if err := json.Unmarshal(pageResp, &page); err != nil {
			log.Printf("github.fetchAPI() -> error unmarshalling response: %v", err)
			continue
		}
		// ...and append.
		result = append(result, page...)
	}

	return result
}

// Makes a single request to the GitHub API.
func executeRequest(url string) (*http.Response, error) {
	gitlabPAT := os.Getenv("GITHUB_PAT")
	if gitlabPAT == "" {
		log.Fatalln("GITHUB_PAT is not set")
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("github.executeRequest() -> error creating the request: %v", err)
	}

	authHeader := fmt.Sprintf("Bearer %s", gitlabPAT)
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2026-03-10")

	// Send request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github.executeRequest() -> error sending the request: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github.executeRequest() -> request status code error: %s with URL: %s", resp.Status, url)
	}

	return resp, nil
}

// Extracts the next link from the paginated GitHub API response.
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
