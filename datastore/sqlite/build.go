package sqlite

import (
	"database/sql"
	"github.com/ogra1/fabrica/domain"
	"github.com/rs/xid"
	"log"
)

const createBuildTableSQL string = `
	CREATE TABLE IF NOT EXISTS build (
		id               varchar(200) primary key not null,
		name             varchar(200) not null,
		repo             varchar(200) not null,
		status           varchar(20) default '',
		created          timestamp default current_timestamp,
        download         varchar(20) default ''
	)
`

const addBuildSQL = `
	INSERT INTO build (id, name, repo) VALUES ($1, $2, $3)
`
const updateBuildSQL = `
	UPDATE build SET status=$1 WHERE id=$2
`
const updateBuildDownloadSQL = `
	UPDATE build SET download=$1 WHERE id=$2
`
const listBuildSQL = `
	SELECT id, name, repo, status, created, download
	FROM build
	ORDER BY created DESC
`
const getBuildSQL = `
	SELECT id, name, repo, status, created, download
	FROM build
	WHERE id=$1
`

// BuildCreate stores a new build request
func (db *DB) BuildCreate(name, repo string) (string, error) {
	id := xid.New()
	_, err := db.Exec(addBuildSQL, id.String(), name, repo)
	return id.String(), err
}

// BuildUpdate updates a build request
func (db *DB) BuildUpdate(id, status string) error {
	_, err := db.Exec(updateBuildSQL, status, id)
	return err
}

// BuildUpdateDownload updates a build request's download file path
func (db *DB) BuildUpdateDownload(id, download string) error {
	_, err := db.Exec(updateBuildDownloadSQL, download, id)
	return err
}

// BuildList get the list of builds
func (db *DB) BuildList() ([]domain.Build, error) {
	logs := []domain.Build{}
	rows, err := db.Query(listBuildSQL)
	if err != nil {
		return logs, err
	}
	defer rows.Close()

	for rows.Next() {
		r := domain.Build{}
		err := rows.Scan(&r.ID, &r.Name, &r.Repo, &r.Status, &r.Created, &r.Download)
		if err != nil {
			return logs, err
		}
		logs = append(logs, r)
	}

	return logs, nil
}

// BuildGet fetches a build with its logs
func (db *DB) BuildGet(id string) (domain.Build, error) {
	r := domain.Build{}
	err := db.QueryRow(getBuildSQL, id).Scan(&r.ID, &r.Name, &r.Repo, &r.Status, &r.Created, &r.Download)
	switch {
	case err == sql.ErrNoRows:
		return r, err
	case err != nil:
		log.Printf("Error retrieving database build: %v\n", err)
		return r, err
	}

	logs, err := db.BuildLogList(id)
	if err != nil {
		return r, err
	}
	r.Logs = logs
	return r, nil
}
