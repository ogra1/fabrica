package datastore

import "github.com/ogra1/fabrica/domain"

// Datastore interface for the database logic
type Datastore interface {
	BuildList() ([]domain.Build, error)
	BuildCreate(name, repo string) (string, error)
	BuildUpdate(id, status string, duration int) error
	BuildUpdateDownload(id, download string) error
	BuildGet(id string) (domain.Build, error)
	BuildDelete(id string) error

	BuildLogCreate(id, message string) error
	BuildLogList(id string) ([]domain.BuildLog, error)

	RepoCreate(name, repo string) (string, error)
	RepoGet(id string) (domain.Repo, error)
	RepoList(watch bool) ([]domain.Repo, error)
	RepoUpdateHash(id, hash string) error
}
