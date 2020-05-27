package writecloser

import (
	"github.com/ogra1/fabrica/datastore"
	"log"
	"strings"
	"sync"
)

// DownloadWriteCloser writes log lines to the database
type DownloadWriteCloser struct {
	lock      sync.RWMutex
	BuildID   string
	Datastore datastore.Datastore
	filename  string
}

// NewDownloadWriteCloser creates a new database write-closer
func NewDownloadWriteCloser(buildID string, ds datastore.Datastore) *DownloadWriteCloser {
	return &DownloadWriteCloser{
		BuildID:   buildID,
		Datastore: ds,
	}
}

// Write writes a log message to the database
func (dwn *DownloadWriteCloser) Write(b []byte) (int, error) {
	dwn.lock.Lock()
	defer dwn.lock.Unlock()

	s := string(b)

	// Check if we have the snapped line
	log.Println(s)
	if strings.Contains(s, "Snapped") {
		parts := strings.Split(s, " ")
		if len(parts) > 0 {
			dwn.filename = dwn.cleanFilename(parts[1])
		}
	}

	if err := dwn.Datastore.BuildLogCreate(dwn.BuildID, s); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (dwn *DownloadWriteCloser) cleanFilename(filePath string) string {
	// Get the snap file name (without the extra text at the end)
	parts := strings.Split(filePath, ".snap")
	return parts[0] + ".snap"
}

// Close is a noop to fulfill the interface
func (dwn *DownloadWriteCloser) Close() error {
	// Noop
	return nil
}

// Filename retrieves the filename of the snap
func (dwn *DownloadWriteCloser) Filename() string {
	dwn.lock.RLock()
	defer dwn.lock.RUnlock()

	return dwn.filename
}
