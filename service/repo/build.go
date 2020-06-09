package repo

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/ogra1/fabrica/datastore"
	"github.com/ogra1/fabrica/domain"
	"github.com/ogra1/fabrica/service"
	"github.com/ogra1/fabrica/service/lxd"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const (
	statusInProgress    = "in-progress"
	statusFailed        = "failed"
	statusComplete      = "complete"
	downloadFileMessage = "Archived snap package: "
)

// BuildSrv interface for building images
type BuildSrv interface {
	Build(repoID string) (string, error)
	List() ([]domain.Build, error)
	BuildGet(id string) (domain.Build, error)
	BuildDelete(id string) error

	RepoCreate(repo, branch, keyID string) (string, error)
	RepoList(watch bool) ([]domain.Repo, error)
	RepoDelete(id string, deleteBuilds bool) error
}

// BuildService implements a build service
type BuildService struct {
	Datastore datastore.Datastore
	LXDSrv    lxd.Service
}

// NewBuildService creates a new build service
func NewBuildService(ds datastore.Datastore, lx lxd.Service) *BuildService {
	return &BuildService{
		Datastore: ds,
		LXDSrv:    lx,
	}
}

// Build starts a build with lxd
func (bld *BuildService) Build(repoID string) (string, error) {
	// Get the repo from the ID
	repo, err := bld.Datastore.RepoGet(repoID)
	if err != nil {
		return "", fmt.Errorf("cannot find the repository: %v", err)
	}

	// Store the build request
	buildID, err := bld.Datastore.BuildCreate(repo.Name, repo.Repo, repo.Branch)
	if err != nil {
		return buildID, fmt.Errorf("error storing build request: %v", err)
	}

	// Start the build in a go routine
	go bld.requestBuild(repo, buildID)

	return buildID, nil
}

func (bld *BuildService) requestBuild(repo domain.Repo, buildID string) error {
	start := time.Now()
	// Update build status
	_ = bld.Datastore.BuildUpdate(buildID, statusInProgress, 0)

	// Clone the repo and get the last commit tag
	repoPath, hash, err := bld.cloneRepo(repo)
	if err != nil {
		log.Println("Cloning repository:", err)
		duration := time.Now().Sub(start).Seconds()
		_ = bld.Datastore.BuildLogCreate(buildID, fmt.Sprintf("Error cloning repo: %v\n", err))
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed, int(duration))
		return err
	}
	log.Printf("Cloned repo: %s (%s)\n", repoPath, hash)
	bld.Datastore.BuildLogCreate(buildID, fmt.Sprintf("Cloned repo: %s (%s)\n", repoPath, hash))

	// Find the snapcraft.yaml file
	f, err := bld.findSnapcraftYAML(repoPath)
	if err != nil {
		log.Println("Find snapcraft.yaml:", err)
		duration := time.Now().Sub(start).Seconds()
		_ = bld.Datastore.BuildLogCreate(buildID, fmt.Sprintf("Find snapcraft.yaml: %v\n", err))
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed, int(duration))
		return err
	}
	log.Printf("snapcraft.yaml: %s\n", f)
	bld.Datastore.BuildLogCreate(buildID, fmt.Sprintf("snapcraft.yaml: %s\n", f))

	// Get the distro from looking at the `base` keyword
	distro, err := bld.getDistroFromYAML(f)
	if err != nil {
		log.Println("Get distro:", err)
		duration := time.Now().Sub(start).Seconds()
		_ = bld.Datastore.BuildLogCreate(buildID, fmt.Sprintf("Get distro: %v\n", err))
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed, int(duration))
		return err
	}
	bld.Datastore.BuildLogCreate(buildID, fmt.Sprintf("Distro: %s\n", distro))

	// Clean up the cloned repo
	_ = os.RemoveAll(repoPath)

	// Run the build in an LXD container
	if err := bld.LXDSrv.RunBuild(buildID, repo.Name, repo.Repo, repo.Branch, distro); err != nil {
		duration := time.Now().Sub(start).Seconds()
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed, int(duration))
		return err
	}

	// Update the repo's last commit
	_ = bld.Datastore.RepoUpdateHash(repo.ID, hash)

	// Mark the build as complete
	duration := time.Now().Sub(start).Seconds()
	_ = bld.Datastore.BuildUpdate(buildID, statusComplete, int(duration))
	return nil
}

// cloneRepo the repo and return the path and tag
func (bld *BuildService) cloneRepo(r domain.Repo) (string, string, error) {
	// Clone the repo
	p := service.GetPath(r.ID)
	log.Println("git", "clone", "--depth", "1", r.Repo, p)
	refBranch := plumbing.NewBranchReferenceName(r.Branch)
	gitRepo, err := git.PlainClone(p, false, &git.CloneOptions{
		URL:           r.Repo,
		ReferenceName: refBranch,
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		log.Println("Error cloning repo:", err)
		return "", "", err
	}

	// Get the last commit hash
	log.Println("git", "ls-remote", "--heads", p)
	ref, err := gitRepo.Head()
	if err != nil {
		return "", "", err
	}
	return p, ref.Hash().String(), nil
}

func (bld *BuildService) findSnapcraftYAML(p string) (string, error) {
	// Check the root directory for snapcraft.yaml
	f := path.Join(p, "snapcraft.yaml")
	log.Println("Checking path:", f)
	_, err := os.Stat(f)
	if err == nil {
		return f, nil
	}

	// Check the root directory for snapcraft.yaml
	f = path.Join(p, "snap", "snapcraft.yaml")
	log.Println("Checking path:", f)
	_, err = os.Stat(f)
	if err == nil {
		return f, nil
	}

	return "", fmt.Errorf("cannot file snapcraft.yaml in repository")
}

func (bld *BuildService) getDistroFromYAML(f string) (string, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return "", err
	}

	keys := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &keys); err != nil {
		return "", err
	}

	base, ok := keys["base"].(string)
	if !ok {
		// Default to xenial when there is no base defined
		return "xenial", nil
	}

	// Convert the base to a distro
	switch base {
	case "core18":
		return "bionic", nil
	case "core20":
		return "focal", nil
	default:
		return "xenial", nil
	}
}

// checkForDownloadFile parses the message to see if we have the download file path
func (bld *BuildService) checkForDownloadFile(buildID, message string) {
	p := strings.TrimPrefix(message, downloadFileMessage)
	if err := bld.Datastore.BuildUpdateDownload(buildID, p); err != nil {
		log.Println("Error storing download path:", err)
	}
}

func nameFromRepo(repo string) string {
	base := path.Base(repo)
	return strings.TrimSuffix(base, ".git")
}
