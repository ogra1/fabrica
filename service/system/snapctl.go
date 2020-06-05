package system

import (
	"log"
	"os/exec"
)

// SnapCtlGet fetches a snap configuration option
func (c *Service) SnapCtlGet(key string) (string, error) {
	out, err := exec.Command("snapctl", "get", key).Output()
	if err != nil {
		return "", err
	}
	log.Println("---snapctl:", key, string(out))
	return string(out), nil
}

// SnapCtlGetBool fetches a snap configuration option that is boolean
func (c *Service) SnapCtlGetBool(key string) bool {
	value, err := c.SnapCtlGet(key)
	if err != nil {
		log.Println("Error calling snapctl:", err)
		return false
	}

	return value == "true"
}
