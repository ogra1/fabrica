package repo

import (
	"fmt"
	"github.com/ogra1/fabrica/domain"
)

// RepoCreate creates a new repo
func (bld *BuildService) RepoCreate(repo string) (string, error) {
	// Store the build request
	name := nameFromRepo(repo)
	repoID, err := bld.Datastore.RepoCreate(name, repo)
	if err != nil {
		return repoID, fmt.Errorf("error storing repo: %v", err)
	}

	return repoID, nil
}

// RepoList returns a list of repos
func (bld *BuildService) RepoList(watch bool) ([]domain.Repo, error) {
	return bld.Datastore.RepoList(watch)
}
