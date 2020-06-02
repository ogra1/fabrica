package system

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
)

// Srv interface for system resources
type Srv interface {
	CPU() (float64, error)
	Memory() (uint64, error)
}

// Service implements a system service
type Service struct {
}

// NewSystemService creates a new system service
func NewSystemService() *Service {
	return &Service{}
}

// CPU returns the current CPU usage
func (c *Service) CPU() (float64, error) {
	vv, err := cpu.Times(false)
	if err != nil {
		log.Printf("Error getting cpu usage: %v\n", err)
		return 0, err
	}

	var total float64
	for _, v := range vv {
		total += v.Total()
	}
	return total, nil
}

// Memory returns the current memory usage
func (c *Service) Memory() (uint64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v\n", err)
		return 0, err
	}

	return v.Total, nil
}
