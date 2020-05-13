package service

import (
	"bufio"
	"fmt"
	"github.com/ogra1/fabrica/datastore"
	"github.com/ogra1/fabrica/domain"
	"io"
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
)

// BuildSrv interface for building images
type BuildSrv interface {
	Build(repo string) (string, error)
	List() ([]domain.Build, error)
	BuildGet(id string) (domain.Build, error)
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
func (bld *BuildService) Build(repo string) (string, error) {
	// Store the build request
	name := nameFromRepo(repo)
	buildID, err := bld.Datastore.BuildCreate(name, repo)
	if err != nil {
		return buildID, fmt.Errorf("error storing build request: %v", err)
	}

	// Start the build in a go routine
	go bld.requestBuild(repo, buildID)

	return buildID, nil
}

func (bld *BuildService) requestBuild(repo string, buildID string) error {
	// Update build status
	_ = bld.Datastore.BuildUpdate(buildID, statusInProgress)

	// Set up the build command
	p := path.Join(os.Getenv("SNAP"), "bin/build.py")

	cmd := exec.Command(p, repo, buildID)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed)
		log.Println(err)
		return err
	}

	// Start the build
	if err := cmd.Start(); err != nil {
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed)
		log.Println(err)
		return err
	}

	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for s.Scan() {
		if err := bld.Datastore.BuildLogCreate(buildID, s.Text()); err != nil {
			_ = bld.Datastore.BuildUpdate(buildID, statusFailed)
			return fmt.Errorf("error storing log: %v", err)
		}

		if strings.HasPrefix(s.Text(), downloadFileMessage) {
			bld.checkForDownloadFile(buildID, s.Text())
		}
	}

	if err := cmd.Wait(); err != nil {
		_ = bld.Datastore.BuildUpdate(buildID, statusFailed)
		log.Println(err)
		return err
	}

	_ = bld.Datastore.BuildUpdate(buildID, statusComplete)
	return nil
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
