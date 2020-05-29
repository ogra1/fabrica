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

// RepoDelete removes a repo and optionally removes its builds
func (bld *BuildService) RepoDelete(id string, deleteBuilds bool) error {
	if !deleteBuilds {
		// Just remove the repo record
		return bld.Datastore.RepoDelete(id)
	}

	// Get the repo name from its ID
	repo, err := bld.Datastore.RepoGet(id)
	if err != nil {
		return fmt.Errorf("cannot find repo: %v", err)
	}

	// Get the builds for this repo
	builds, err := bld.Datastore.BuildListForRepo(repo.Name)
	if err != nil {
		return fmt.Errorf("error finding builds: %v", err)
	}

	// Remove each of the repo's builds
	for _, b := range builds {
		if err := bld.BuildDelete(b.ID); err != nil {
			return fmt.Errorf("error removing build (%s): %v", b.ID, err)
		}
	}

	// Remove the repo record
	return bld.Datastore.RepoDelete(id)
}
