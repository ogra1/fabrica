package repo

import (
	"bufio"
	"fmt"
	"github.com/ogra1/fabrica/datastore"
	"github.com/ogra1/fabrica/domain"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	statusInProgress    = "in-progress"
	statusFailed        = "failed"
	statusComplete      = "complete"
	downloadFileMessage = "Archived snap package: "
	snapData            = "SNAP_DATA"
)

// BuildSrv interface for building images
type BuildSrv interface {
	Build(repo string) (string, error)
	List() ([]domain.Build, error)
	BuildGet(id string) (domain.Build, error)
	RepoCreate(repo string) (string, error)
	RepoList(watch bool) ([]domain.Repo, error)
}

// BuildService implements a build service
type BuildService struct {
	Datastore datastore.Datastore
}

// NewBuildService creates a new build service
func NewBuildService(ds datastore.Datastore) *BuildService {
	return &BuildService{
		Datastore: ds,
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
	buildID, err := bld.Datastore.BuildCreate(repo.Name, repo.Repo)
	if err != nil {
		return buildID, fmt.Errorf("error storing build request: %v", err)
	}

	// Start the build in a go routine
	go bld.requestBuild(repo, buildID)

	return buildID, nil
}

func (bld *BuildService) requestBuild(repo domain.Repo, buildID string) error {
	// Update build status
	_ = bld.Datastore.BuildUpdate(buildID, statusInProgress)

	// Clone the repo and get the last commit tag
	repoPath, hash, err := bld.cloneRepo(repo)
	if err != nil {
		log.Println("Cloning repository:", err)
		return err
	}
	log.Printf("Cloned repo: %s (%s)\n", repoPath, hash)
	bld.Datastore.BuildLogCreate(buildID, fmt.Sprintf("Cloned repo: %s (%s)\n", repoPath, hash))

	// Find the snapcraft.yaml file
	f, err := bld.findSnapcraftYAML(repoPath)
	if err != nil {
		log.Println("Find snapcraft.yaml:", err)
		return err
	}
	log.Printf("snapcraft.yaml: %s\n", f)
	bld.Datastore.BuildLogCreate(buildID, fmt.Sprintf("snapcraft.yaml: %s\n", f))

	// Get the distro from looking at the `base` keyword
	distro, err := bld.getDistroFromYAML(f)
	if err != nil {
		log.Println("Get distro:", err)
		return err
	}
	bld.Datastore.BuildLogCreate(buildID, fmt.Sprintf("Distro: %s\n", f))

	// Run the build via the python script
	cmd, err := bld.runBuild(repo, buildID, distro, err)
	if err != nil {
		log.Println("Run build:", err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed)
		log.Println(err)
		return err
	}

	// Update the repo's last commit
	_ = bld.Datastore.RepoUpdateHash(repo.ID, hash)

	// Mark the build as complete
	_ = bld.Datastore.BuildUpdate(buildID, statusComplete)
	return nil
}

func (bld *BuildService) runBuild(repo domain.Repo, buildID string, distro string, err error) (*exec.Cmd, error) {
	// Set up the build command
	p := path.Join(os.Getenv("SNAP"), "bin/build.py")

	cmd := exec.Command(p, repo.Repo, buildID, distro)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed)
		log.Println(err)
		return nil, err
	}

	// Start the build
	if err := cmd.Start(); err != nil {
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed)
		log.Println(err)
		return nil, err
	}

	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for s.Scan() {
		if err := bld.Datastore.BuildLogCreate(buildID, s.Text()); err != nil {
			_ = bld.Datastore.BuildUpdate(buildID, statusFailed)
			return nil, fmt.Errorf("error storing log: %v", err)
		}

		if strings.HasPrefix(s.Text(), downloadFileMessage) {
			bld.checkForDownloadFile(buildID, s.Text())
		}
	}
	return cmd, nil
}

// cloneRepo the repo and return the path and tag
func (bld *BuildService) cloneRepo(r domain.Repo) (string, string, error) {
	// Clone the repo
	p := getPath(r.ID)
	log.Println("git", "clone", "--depth", "1", r.Repo, p)
	o, err := exec.Command("git", "clone", "--depth", "1", r.Repo, p).Output()
	if err != nil {
		log.Println("Clone:", string(o))
		return "", "", err
	}

	// Get the last commit hash
	out, err := exec.Command("git", "ls-remote", "--heads", p).Output()
	if err != nil {
		return "", "", err
	}
	refs := strings.Split(string(out), " ")
	return p, refs[0], nil
}

func (bld *BuildService) findSnapcraftYAML(p string) (string, error) {
	// Check the root directory for snapcraft.yaml
	f := path.Join(p, "snapcraft.yaml")
	_, err := os.Stat(f)
	if os.IsExist(err) {
		return f, nil
	}

	// Check the root directory for snapcraft.yaml
	f = path.Join(p, "snap", "snapcraft.yaml")
	_, err = os.Stat(f)
	if os.IsExist(err) {
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
	if err := yaml.Unmarshal(data, data); err != nil {
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

// List returns a list of the builds that have been requested
func (bld *BuildService) List() ([]domain.Build, error) {
	return bld.Datastore.BuildList()
}

// BuildGet returns a build with its logs
func (bld *BuildService) BuildGet(id string) (domain.Build, error) {
	return bld.Datastore.BuildGet(id)
}

func nameFromRepo(repo string) string {
	base := path.Base(repo)
	return strings.TrimSuffix(base, ".git")
}

func getPath(p string) string {
	return path.Join(os.Getenv(snapData), p)
}
