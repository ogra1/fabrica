package writecloser

import (
	"github.com/ogra1/fabrica/datastore"
	"log"
	"strings"
	"sync"
)

// FlagWriteCloser is a writer-closer that looks for a specific message
type FlagWriteCloser struct {
	lock      sync.RWMutex
	Contains  string
	Datastore datastore.Datastore
	found     bool
}

// NewFlagWriteCloser creates a new flag write-closer
func NewFlagWriteCloser(contains string) *FlagWriteCloser {
	return &FlagWriteCloser{
		Contains: contains,
	}
}

// Write writes a log message to the database
func (wc *FlagWriteCloser) Write(b []byte) (int, error) {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	log.Println(string(b))

	if strings.Contains(string(b), wc.Contains) {
		wc.found = true
	}
	return len(b), nil
}

// Close is a noop to fulfill the interface
func (wc *FlagWriteCloser) Close() error {
	// Noop
	return nil
}

// Found identifies if the string has been found
func (wc *FlagWriteCloser) Found() bool {
	wc.lock.RLock()
	defer wc.lock.RUnlock()

	return wc.found
}
