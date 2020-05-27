package service

import (
	"os"
	"path"
)

const (
	snapCommon = "SNAP_COMMON"
)

// GetPath gets a path from SNAP_COMMON
func GetPath(p string) string {
	return path.Join(os.Getenv(snapCommon), p)
}
