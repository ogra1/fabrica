package lxd

import (
	"github.com/ogra1/fabrica/domain"
	"os/exec"
	"syscall"
)

var plugs = []string{"lxd", "mount-observe", "system-observe", "ssh-keys"}

// GetImageAlias checks of an image alias is available
func (lx *LXD) GetImageAlias(name string) error {
	// Connect to the lxd service
	c, err := lx.connect()
	if err != nil {
		return err
	}

	// Check if the alias exists (the image could still be loading)
	_, _, err = c.GetImageAlias(name)
	return err
}

// CheckConnections checks the snap interfaces are connected
func (lx *LXD) CheckConnections() []domain.SettingAvailable {
	results := []domain.SettingAvailable{}

	for _, p := range plugs {
		exitCode := runCommand("snapctl", "is-connected", p)

		// Store the setting
		results = append(results, domain.SettingAvailable{
			Name:      p,
			Available: exitCode == 0,
		})
	}
	return results
}

func runCommand(name string, args ...string) int {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			return ws.ExitStatus()
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		return ws.ExitStatus()
	}
	return 1
}
