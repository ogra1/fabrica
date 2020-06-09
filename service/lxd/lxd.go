package lxd

import (
	"fmt"
	lxd "github.com/lxc/lxd/client"
	"github.com/ogra1/fabrica/datastore"
	"github.com/ogra1/fabrica/domain"
	"github.com/ogra1/fabrica/service/system"
	"log"
	"os"
	"path"
	"time"
)

// Service is the interface for the LXD service
type Service interface {
	RunBuild(buildID, name, repo, branch, keyID, distro string) error
	GetImageAlias(name string) error
	CheckConnections() []domain.SettingAvailable
	StopAndDeleteContainer(name string) error
}

// LXD services
type LXD struct {
	Datastore datastore.Datastore
	SystemSrv system.Srv
}

// NewLXD creates a new LXD client
func NewLXD(ds datastore.Datastore, sysSrv system.Srv) *LXD {
	return &LXD{
		Datastore: ds,
		SystemSrv: sysSrv,
	}
}

// connect opens a connection to the LXD service
func (lx *LXD) connect() (lxd.InstanceServer, error) {
	// Get LXD socket path
	lxdSocket, err := lxdSocketPath()
	if err != nil {
		log.Println("Error with lxd socket:", err)
		return nil, err
	}

	// Connect to LXD over the Unix socket
	c, err := lxd.ConnectLXDUnix(lxdSocket, nil)
	if err != nil {
		log.Println("Error connecting to LXD:", err)
		return nil, err
	}
	return c, nil
}

// RunBuild runs a build by launching a container for the build request
func (lx *LXD) RunBuild(buildID, name, repo, branch, keyID, distro string) error {
	// Create a new LXD connection
	c, err := lx.connect()
	if err != nil {
		lx.Datastore.BuildLogCreate(buildID, err.Error())
		return err
	}

	// Run the build
	run := newRunner(buildID, lx.Datastore, lx.SystemSrv, c)
	return run.runBuild(name, repo, branch, keyID, distro)
}

// StopAndDeleteContainer stops and removes a container
func (lx *LXD) StopAndDeleteContainer(name string) error {
	c, err := lx.connect()
	if err != nil {
		log.Println("Error connecting to LXD:", err)
		return err
	}

	return stopAndDeleteContainer(c, name)
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
