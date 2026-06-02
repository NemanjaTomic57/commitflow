package gitlab

import "time"

// Response type for the API
type responseType interface {
	commit | project
}

type project struct {
	ID                     int              `json:"id"`
	Description            string           `json:"description"`
	PathWithNamespace      string           `json:"path_with_namespace"`
	CreatedAt              string           `json:"created_at"`
	WebURL                 string           `json:"web_url"`
	GitlabProjectNamespace projectNamespace `json:"namespace"`
}

type projectNamespace struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
	FullPath string `json:"full_path"`
	WebURL   string `json:"web_url"`
}

type commit struct {
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
