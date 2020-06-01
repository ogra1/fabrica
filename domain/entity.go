package domain

import "time"

// Repo is the requested repository to watch
type Repo struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Repo       string    `json:"repo"`
	Branch     string    `json:"branch"`
	LastCommit string    `json:"hash"`
	Created    time.Time `json:"created"`
	Modified   time.Time `json:"modified"`
}

// Build is the requested build
type Build struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Repo     string     `json:"repo"`
	Branch   string     `json:"branch"`
	Status   string     `json:"status"`
	Download string     `json:"download"`
	Duration int        `json:"duration"`
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

// SettingAvailable is a generic response for a setting
type SettingAvailable struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
}
