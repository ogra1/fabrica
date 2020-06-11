package repo

import (
	"fmt"
	"github.com/ogra1/fabrica/domain"
	"log"
	"strings"
)

// RepoCreate creates a new repo
func (bld *BuildService) RepoCreate(repo, branch, keyID string) (string, error) {
	// Handle Launchpad git URLs
	if strings.HasPrefix(repo, "git+ssh") {
		repo = strings.Replace(repo, "git+ssh", "ssh", 1)
	}

	// Store the build request
	name := nameFromRepo(repo)
	repoID, err := bld.Datastore.RepoCreate(name, repo, branch, keyID)
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
	log.Println("Repo Delete:", id, deleteBuilds)

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
	builds, err := bld.Datastore.BuildListForRepo(repo.Repo, repo.Branch)
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
