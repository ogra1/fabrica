package watch

import (
	"github.com/ogra1/fabrica/datastore"
	"github.com/ogra1/fabrica/domain"
	"github.com/ogra1/fabrica/service/repo"
	"log"
	"os/exec"
	"strings"
	"time"
)

const (
	tickInterval = 3000
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
			// check for an update
			hash, update, err := srv.checkForUpdates(r)
			if err != nil {
				log.Println("Error checking repository:", err)
				break
			}
			if !update {
				break
			}

			// trigger the build
			if _, err := srv.BuildSrv.Build(r.ID); err != nil {
				log.Println("Error building snap:", err)
				break
			}

			// update the last commit hash
			srv.Datastore.RepoUpdateHash(r.ID, hash)
		}
	}
	ticker.Stop()
}

func (srv *Service) checkForUpdates(r domain.Repo) (string, bool, error) {
	// Get the last commit hash
	out, err := exec.Command("git", "ls-remote", "--heads", r.Repo).Output()
	if err != nil {
		return "", false, err
	}
	refs := strings.Split(string(out), " ")

	return refs[0], r.LastCommit != refs[0], nil
}
