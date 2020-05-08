package domain

import "time"

// Build is the requested build
type Build struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Repo    string     `json:"repo"`
	Status  string     `json:"status"`
	Created time.Time  `json:"created"`
	Logs    []BuildLog `json:"logs,omitempty"`
}

// BuildLog is the log for a build
type BuildLog struct {
	ID      string    `json:"id"`
	BuildID string    `json:"buildId"`
	Message string    `json:"message"`
	Created time.Time `json:"created"`
}
