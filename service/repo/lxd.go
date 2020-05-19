package repo

import (
	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/ogra1/fabrica/datastore"
	"io"
	"log"
	"time"
)

var containerEnv = map[string]string{
	"FLASH_KERNEL_SKIP":           "true",
	"DEBIAN_FRONTEND":             "noninteractive",
	"TERM":                        "xterm",
	"SNAPCRAFT_BUILD_ENVIRONMENT": "host",
}
var containerCmd = [][]string{
	{"apt", "update"},
	{"apt", "-y", "upgrade"},
	{"apt", "-y", "install", "build-essential"},
	{"apt", "-y", "clean"},
	{"snap", "install", "snapcraft", "--classic"},
}

// LXD services
type LXD struct {
	BuildID   string
	Datastore datastore.Datastore
}

// NewLXD creates a new LXD client
func NewLXD(buildID string, ds datastore.Datastore) *LXD {
	return &LXD{
		BuildID:   buildID,
		Datastore: ds,
	}
}

// RunBuild launches an LXD container to start the build
func (lx *LXD) RunBuild(name, repo, distro string) error {
	// Connect to LXD over the Unix socket
	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		return err
	}

	// Generate the container name
	cname := containerName(name)

	// Create and start the LXD
	if err := lx.createAndStartContainer(c, cname, distro); err != nil {
		log.Println("Error creating/starting container:", err)
		lx.Datastore.BuildLogCreate(lx.BuildID, err.Error())
		return err
	}

	// Set up the database writer for the logs
	dbWC := NewDBWriteCloser(lx.BuildID, lx.Datastore)

	// Set up the container
	for _, cmd := range containerCmd {
		if err := lx.runInContainer(c, cname, cmd, dbWC); err != nil {
			log.Println("Command error:", cmd, err)
			lx.Datastore.BuildLogCreate(lx.BuildID, err.Error())
			return err
		}
	}

	// Run the build
	commands := [][]string{
		{"git", "clone", "--progress", repo},
		{"sh", "-cv", "cd /root/" + name + ";", "snapcraft"},
	}
	for _, cmd := range commands {
		if err := lx.runInContainer(c, cname, cmd, dbWC); err != nil {
			log.Println("Command error:", cmd, err)
			lx.Datastore.BuildLogCreate(lx.BuildID, err.Error())
			return err
		}
	}

	// Pull the snap from the container

	// Remove the container
	return nil
}

func (lx *LXD) createAndStartContainer(c lxd.InstanceServer, cname, distro string) error {
	// Container creation request
	req := api.ContainersPost{
		Name: cname,
		Source: api.ContainerSource{
			Type:  "image",
			Alias: distro,
		},
	}

	// Get LXD to create the container (background operation)
	op, err := c.CreateContainer(req)
	if err != nil {
		return err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return err
	}

	// Get LXD to start the container (background operation)
	reqState := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err = c.UpdateContainerState(cname, reqState, "")
	if err != nil {
		return err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (lx *LXD) runInContainer(c lxd.InstanceServer, cname string, command []string, stdOutErr io.WriteCloser) error {
	// Setup the exec request
	req := api.ContainerExecPost{
		Environment: containerEnv,
		Command:     command,
		WaitForWS:   true,
	}

	// Setup the exec arguments (fds)
	args := lxd.ContainerExecArgs{
		Stdout: stdOutErr,
		Stderr: stdOutErr,
	}

	// Get the current state
	op, err := c.ExecContainer(cname, req, &args)
	if err != nil {
		return err
	}

	// Wait for it to complete
	return op.Wait()
}

func containerName(name string) string {
	return name + "-" + string(time.Now().Unix())
}
