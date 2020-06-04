package system

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"log"
	"os"
	"path"
)

// Srv interface for system resources
type Srv interface {
	CPU() (float64, error)
	Memory() (float64, error)
	Disk() (float64, error)
}

// Service implements a system service
type Service struct {
}

// NewSystemService creates a new system service
func NewSystemService() *Service {
	return &Service{}
}

const (
        snapData = "SNAP_DATA"
)

// CPU returns the current CPU usage
func (c *Service) CPU() (float64, error) {
	vv, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting cpu usage: %v\n", err)
		return 0, err
	}

	var total float64
	if len(vv) > 0 {
		total = vv[0]
	}

	return total, nil
}

// Memory returns the current memory usage
func (c *Service) Memory() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v\n", err)
		return 0, err
	}

	return v.UsedPercent, nil
}

// Disk returns the current disk usage
func (c *Service) Disk() (float64, error) {
	// Check the disk space of the host FS not the snap
	v, err := disk.Usage(path.Join(os.Getenv(snapData)))
	if err != nil {
		log.Printf("Error getting disk usage: %v\n", err)
		return 0, err
	}

	return v.UsedPercent, nil
}
