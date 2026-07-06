package client

import "time"

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

type Actor struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
	GravatarID   string `json:"gravatar_id"`
	URL          string `json:"url"`
	AvatarURL    string `json:"avatar_url"`
}

type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Payload struct {
	Action       string       `json:"action,omitempty"`
	Ref          string       `json:"ref,omitempty"`
	RefType      string       `json:"ref_type,omitempty"`
	Head         string       `json:"head,omitempty"`
	Before       string       `json:"before,omitempty"`
	Size         int          `json:"size,omitempty"`
	PushID       int64        `json:"push_id,omitempty"`
	RepositoryID int64        `json:"repository_id,omitempty"`
	DistinctSize int          `json:"distinct_size,omitempty"`
	Commits      []Commit     `json:"commits,omitempty"`
	Issue        *Issue       `json:"issue,omitempty"`
	PullRequest  *PullRequest `json:"pull_request,omitempty"`
	Forkee       *Repo        `json:"forkee,omitempty"`
	Member       *Member      `json:"member,omitempty"`
}

type Commit struct {
	SHA      string `json:"sha"`
	Message  string `json:"message"`
	Distinct bool   `json:"distinct"`
	URL      string `json:"url"`
}

type Issue struct {
	ID     int    `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	URL    string `json:"url"`
}

type Member struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type PullRequest struct {
	ID     int    `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	URL    string `json:"url"`
}
