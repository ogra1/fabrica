package watch

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/ogra1/fabrica/datastore"
	"github.com/ogra1/fabrica/domain"
	"github.com/ogra1/fabrica/service/repo"
	"log"
	"time"
)

const (
	tickInterval = 300
)

// Srv interface for watching repos
type Srv interface {
	Watch()
}

// Service implements a build service
type Service struct {
	BuildSrv  repo.BuildSrv
	Datastore datastore.Datastore
}

// NewWatchService creates a new watch service
func NewWatchService(ds datastore.Datastore, buildSrv repo.BuildSrv) *Service {
	return &Service{
		Datastore: ds,
		BuildSrv:  buildSrv,
	}
}

// Watch service to watch repo updates
// The service runs on an interval and picks up the list of repos from the database
// (fetching the most stale first). It goes through each repo, checking the latest
// commit hash with the one stored in the database. As soon as it identifies one
// repo that needs to be updated, it clones the repo and parse the snapcraft.yaml
// to get the base. Then the lxd build process is started for that repo and the
// commit hash is updated in the database. The watch service then has complete its
// cycle and will wait for the next interval. So only one repo is processed in
// each cycle.
func (srv *Service) Watch() {
	// On an interval...
	ticker := time.NewTicker(time.Second * tickInterval)
	for range ticker.C {
		// Get the repo list
		records, err := srv.BuildSrv.RepoList(true)
		if err != nil {
			log.Println("Error fetching repositories:", err)
			break
		}

		for _, r := range records {
			log.Println("Check repo:", r.Repo)
			// check for an update
			hash, update, err := srv.checkForUpdates(r)
			if err != nil {
				log.Println("Error checking repository:", err)
				// check the next repo
				continue
			}
			if !update {
				// no update so check the next repo
				continue
			}

			// update the last commit hash (to avoid repeating builds)
			srv.Datastore.RepoUpdateHash(r.ID, hash)

			// trigger the build
			if _, err := srv.BuildSrv.Build(r.ID); err != nil {
				log.Println("Error building snap:", err)
				// check the next repo
				continue
			}

			// Don't process any more repos until the next cycle
			break
		}
	}
	ticker.Stop()
}

func (srv *Service) checkForUpdates(r domain.Repo) (string, bool, error) {
	// Get the last commit hash
	log.Println("git", "ls-remote", "--heads", r.Repo)
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{r.Repo},
	})

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return "", false, err
	}

	for _, ref := range refs {
		if checkBranch(r.Branch, ref) {
			return ref.Hash().String(), r.LastCommit != ref.Hash().String(), nil
		}
	}
	return "", false, fmt.Errorf("cannot find the repo HEAD")
}

func checkBranch(branch string, ref *plumbing.Reference) bool {
	name := fmt.Sprintf("refs/heads/%s", branch)
	return ref.Name().IsBranch() && ref.Name().String() == name
}
