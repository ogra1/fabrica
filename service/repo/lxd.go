package repo

import (
	"fmt"
	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/ogra1/fabrica/datastore"
	"io"
	"log"
	"os"
	"path"
	"strings"
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
	{"snap", "list"},
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
	log.Println("Run build:", name, repo, distro)
	log.Println("Creating and starting container")
	lx.Datastore.BuildLogCreate(lx.BuildID, "milestone: Creating and starting container")

	// Get LXD socket path
	lxdSocket, err := lxdSocketPath()
	if err != nil {
		log.Println("Error with lxd socket:", err)
		lx.Datastore.BuildLogCreate(lx.BuildID, err.Error())
		return err
	}

	// Connect to LXD over the Unix socket
	c, err := lxd.ConnectLXDUnix(lxdSocket, nil)
	if err != nil {
		log.Println("Error creating/starting container:", err)
		lx.Datastore.BuildLogCreate(lx.BuildID, err.Error())
		return err
	}

	// Generate the container name
	cname := containerName(name)

	// Create and start the LXD
	lx.Datastore.BuildLogCreate(lx.BuildID, "milestone: Create container "+cname)
	if err := lx.createAndStartContainer(c, cname, distro); err != nil {
		log.Println("Error creating/starting container:", err)
		lx.Datastore.BuildLogCreate(lx.BuildID, err.Error())
		return err
	}

	// Set up the database writer for the logs
	dbWC := NewDBWriteCloser(lx.BuildID, lx.Datastore)

	// Wait for the network to be running
	lx.Datastore.BuildLogCreate(lx.BuildID, "milestone: Waiting for the network")
	lx.waitForNetwork(c, cname)
	lx.Datastore.BuildLogCreate(lx.BuildID, "milestone: Network is ready")

	// Set up the container
	lx.Datastore.BuildLogCreate(lx.BuildID, "milestone: Install dependencies")
	for _, cmd := range containerCmd {
		lx.Datastore.BuildLogCreate(lx.BuildID, "milestone: "+strings.Join(cmd, " "))
		if err := lx.runInContainer(c, cname, cmd, dbWC); err != nil {
			log.Println("Command error:", cmd, err)
			lx.Datastore.BuildLogCreate(lx.BuildID, err.Error())
			return err
		}
	}

	// Run the build
	commands := [][]string{
		{"git", "clone", "--progress", repo},
		{"cd", "/root/" + name},
		{"snapcraft"},
	}
	lx.Datastore.BuildLogCreate(lx.BuildID, "milestone: Clone repo and build snap")
	for _, cmd := range commands {
		lx.Datastore.BuildLogCreate(lx.BuildID, "milestone: "+strings.Join(cmd, " "))
		if err := lx.runInContainer(c, cname, cmd, dbWC); err != nil {
			log.Println("Command error:", cmd, err)
			lx.Datastore.BuildLogCreate(lx.BuildID, err.Error())
			return err
		}
	}

	// Pull the snap from the container
	lx.Datastore.BuildLogCreate(lx.BuildID, "milestone: Transfer the snap from the container")
	cmd := []string{
		"sh", "-cv", "'ls /root/" + name + "/*.snap'",
	}
	if err := lx.runInContainer(c, cname, cmd, dbWC); err != nil {
		log.Println("Command error:", cmd, err)
		lx.Datastore.BuildLogCreate(lx.BuildID, err.Error())
		return err
	}

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

func (lx *LXD) waitForNetwork(c lxd.InstanceServer, cname string) {
	// Set up the writer to check for a message
	wc := NewFlagWriteCloser("PING")

	// Run a command in the container
	log.Println("Waiting for network...")
	cmd := []string{"ping", "-c1", "8.8.8.8"}
	for {
		_ = lx.runInContainer(c, cname, cmd, wc)
		if wc.Found() {
			break
		}
	}
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
	return fmt.Sprintf("%s-%d", name, time.Now().Unix())
}

// lxdSocketPath finds the socket path for LXD
func lxdSocketPath() (string, error) {
	var ff = []string{
		path.Join(os.Getenv("LXD_DIR"), "unix.socket"),
		"/var/snap/lxd/common/lxd/unix.socket",
		"/var/lib/lxd/unix.socket",
	}

	for _, f := range ff {
		if _, err := os.Stat(f); err == nil {
			return f, nil
		}
	}

	return "", fmt.Errorf("cannot find the LXD socket file")
}
