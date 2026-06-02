package github

import "time"

type responseType interface {
	repository | commitResponse
}

type repository struct {
	ID                        int64     `json:"id"`
	NodeID                    string    `json:"node_id"`
	Name                      string    `json:"name"`
	FullName                  string    `json:"full_name"`
	Private                   bool      `json:"private"`
	Owner                     user      `json:"owner"`
	HTMLURL                   string    `json:"html_url"`
	Description               *string   `json:"description"`
	Fork                      bool      `json:"fork"`
	URL                       string    `json:"url"`
	ForksURL                  string    `json:"forks_url"`
	KeysURL                   string    `json:"keys_url"`
	CollaboratorsURL          string    `json:"collaborators_url"`
	TeamsURL                  string    `json:"teams_url"`
	HooksURL                  string    `json:"hooks_url"`
	IssueEventsURL            string    `json:"issue_events_url"`
	EventsURL                 string    `json:"events_url"`
	AssigneesURL              string    `json:"assignees_url"`
	BranchesURL               string    `json:"branches_url"`
	TagsURL                   string    `json:"tags_url"`
	BlobsURL                  string    `json:"blobs_url"`
	GitTagsURL                string    `json:"git_tags_url"`
	GitRefsURL                string    `json:"git_refs_url"`
	TreesURL                  string    `json:"trees_url"`
	StatusesURL               string    `json:"statuses_url"`
	LanguagesURL              string    `json:"languages_url"`
	StargazersURL             string    `json:"stargazers_url"`
	ContributorsURL           string    `json:"contributors_url"`
	SubscribersURL            string    `json:"subscribers_url"`
	SubscriptionURL           string    `json:"subscription_url"`
	CommitsURL                string    `json:"commits_url"`
	GitCommitsURL             string    `json:"git_commits_url"`
	CommentsURL               string    `json:"comments_url"`
	IssueCommentURL           string    `json:"issue_comment_url"`
	ContentsURL               string    `json:"contents_url"`
	CompareURL                string    `json:"compare_url"`
	MergesURL                 string    `json:"merges_url"`
	ArchiveURL                string    `json:"archive_url"`
	DownloadsURL              string    `json:"downloads_url"`
	IssuesURL                 string    `json:"issues_url"`
	PullsURL                  string    `json:"pulls_url"`
	MilestonesURL             string    `json:"milestones_url"`
	NotificationsURL          string    `json:"notifications_url"`
	LabelsURL                 string    `json:"labels_url"`
	ReleasesURL               string    `json:"releases_url"`
	DeploymentsURL            string    `json:"deployments_url"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	PushedAt                  time.Time `json:"pushed_at"`
	GitURL                    string    `json:"git_url"`
	SSHURL                    string    `json:"ssh_url"`
	CloneURL                  string    `json:"clone_url"`
	SVNURL                    string    `json:"svn_url"`
	Homepage                  *string   `json:"homepage"`
	Size                      int       `json:"size"`
	StargazersCount           int       `json:"stargazers_count"`
	WatchersCount             int       `json:"watchers_count"`
	Language                  *string   `json:"language"`
	HasIssues                 bool      `json:"has_issues"`
	HasProjects               bool      `json:"has_projects"`
	HasDownloads              bool      `json:"has_downloads"`
	HasWiki                   bool      `json:"has_wiki"`
	HasPages                  bool      `json:"has_pages"`
	HasDiscussions            bool      `json:"has_discussions"`
	ForksCount                int       `json:"forks_count"`
	MirrorURL                 *string   `json:"mirror_url"`
	Archived                  bool      `json:"archived"`
	Disabled                  bool      `json:"disabled"`
	OpenIssuesCount           int       `json:"open_issues_count"`
	License                   *license  `json:"license"`
	AllowForking              bool      `json:"allow_forking"`
	IsTemplate                bool      `json:"is_template"`
	WebCommitSignoffRequired  bool      `json:"web_commit_signoff_required"`
	HasPullRequests           bool      `json:"has_pull_requests"`
	PullRequestCreationPolicy string    `json:"pull_request_creation_policy"`
	Topics                    []string  `json:"topics"`
	Visibility                string    `json:"visibility"`
	Forks                     int       `json:"forks"`
	OpenIssues                int       `json:"open_issues"`
	Watchers                  int       `json:"watchers"`
	DefaultBranch             string    `json:"default_branch"`
}

type license struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}

type commitResponse struct {
	SHA         string      `json:"sha"`
	NodeID      string      `json:"node_id"`
	Commit      commit      `json:"commit"`
	URL         string      `json:"url"`
	HTMLURL     string      `json:"html_url"`
	CommentsURL string      `json:"comments_url"`
	Author      *user       `json:"author"`
	Committer   *user       `json:"committer"`
	Parents     []parentRef `json:"parents"`
}

type commit struct {
	Author       commitAuthor `json:"author"`
	Committer    commitAuthor `json:"committer"`
	Message      string       `json:"message"`
	Tree         tree         `json:"tree"`
	URL          string       `json:"url"`
	CommentCount int          `json:"comment_count"`
	Verification verification `json:"verification"`
}

type commitAuthor struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type tree struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type verification struct {
	Verified   bool    `json:"verified"`
	Reason     string  `json:"reason"`
	Signature  *string `json:"signature"`
	Payload    *string `json:"payload"`
	VerifiedAt *string `json:"verified_at"`
}

type user struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type parentRef struct {
	SHA     string `json:"sha"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}
