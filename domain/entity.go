package domain

import "time"

// Repo is the requested repository to watch
type Repo struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Repo     string    `json:"repo"`
	LastHash string    `json:"hash"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// Build is the requested build
type Build struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Repo     string     `json:"repo"`
	Status   string     `json:"status"`
	Download string     `json:"download"`
	Created  time.Time  `json:"created"`
	Logs     []BuildLog `json:"logs,omitempty"`
}

// BuildLog is the log for a build
type BuildLog struct {
	ID      string    `json:"id"`
	BuildID string    `json:"buildId"`
	Message string    `json:"message"`
	Created time.Time `json:"created"`
}
