package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
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

type GitlabAPIResponse interface {
	GitlabProject | GitlabCommit
}

var baseURL = "https://gitlab.com/api/v4"

func fetchGitlabAPI(url string) *http.Response {
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

func extractBodyFromResponse(resp *http.Response) []byte {
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("extractBodyFromResponse() -> error reading response body:", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("extractBodyFromResponse() -> request status code error: %s\n%s", resp.Status, string(body))
	}

	return body
}

func getNextLink(resp *http.Response) string {
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

func produceKafkaEvents[T GitlabAPIResponse](p *kafka.Producer, resp []byte, topic string) {
	var object []T

	err := json.Unmarshal(resp, &object)
	if err != nil {
		log.Fatal("produceKafkaEvents() -> error unmarshalling JSON:", err)
	}

	// Get results back from producing to Kafka and print to console
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("produceKafkaEvents() -> delivery failed %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("produceKafkaEvents() -> delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce to Kafka topic
	for _, project := range object {
		projectBytes, err := json.Marshal(project)
		if err != nil {
			log.Println("produceKafkaEvents() -> error marshalling project:", err)
			continue
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          projectBytes,
		}, nil)
		if err != nil {
			log.Println("produceKafkaEvents() -> error producing message:", err)
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("main() -> error loading .env file")
	}

	bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		log.Fatal("KAFKA_BOOTSTRAP_SERVER is not set")
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "192.168.56.10:9092"})
	if err != nil {
		log.Fatal("main() -> error creating Kafka producer: ", err)
	}
	defer p.Close()

	// Fetch Gitlab project IDs
	projectIDs := fetchGitlabProjectIDs()

	for _, id := range projectIDs {
		url := baseURL + fmt.Sprintf("/projects/%d/repository/commits", id)
		topic := "git.commits"

		for url != "" {
			resp := fetchGitlabAPI(url)
			url = getNextLink(resp)
			body := extractBodyFromResponse(resp)
			resp.Body.Close()
			produceKafkaEvents[GitlabCommit](p, body, topic)
		}

		p.Flush(15 * 1000)
	}
}

func fetchGitlabProjectIDs() []int {
	var projects []GitlabProject
	var projectIDs []int
	url := baseURL + "/projects?owned=true"

	for url != "" {
		resp := fetchGitlabAPI(url)
		url = getNextLink(resp)
		body := extractBodyFromResponse(resp)
		resp.Body.Close()

		json.Unmarshal(body, &projects)
		for _, project := range projects {
			projectIDs = append(projectIDs, project.ID)
		}
	}

	return projectIDs
}
