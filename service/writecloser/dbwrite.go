package writecloser

import "github.com/ogra1/fabrica/datastore"

// DBWriteCloser writes log lines to the database
type DBWriteCloser struct {
	BuildID   string
	Datastore datastore.Datastore
}

// NewDBWriteCloser creates a new database write-closer
func NewDBWriteCloser(buildID string, ds datastore.Datastore) *DBWriteCloser {
	return &DBWriteCloser{
		BuildID:   buildID,
		Datastore: ds,
	}
}

// Write writes a log message to the database
func (dbw *DBWriteCloser) Write(b []byte) (int, error) {
	if err := dbw.Datastore.BuildLogCreate(dbw.BuildID, string(b)); err != nil {
		return 0, err
	}
	return len(b), nil
}

// Close is a noop to fulfill the interface
func (dbw *DBWriteCloser) Close() error {
	// Noop
	return nil
}
