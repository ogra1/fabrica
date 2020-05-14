package sqlite

import (
	"database/sql"
	"github.com/ogra1/fabrica/domain"
	"github.com/rs/xid"
	"log"
)

const createRepoTableSQL string = `
	CREATE TABLE IF NOT EXISTS repo (
		id               varchar(200) primary key not null,
		name             varchar(200) not null,
		location         varchar(200) UNIQUE not null,
		hash             varchar(200) default '',
		created          timestamp default current_timestamp,
		modified         timestamp default current_timestamp
	)
`

const addRepoSQL = `
	INSERT INTO repo (id, name, location) VALUES ($1, $2, $3)
`
const listRepoSQL = `
	SELECT id, name, location, hash, created, modified
	FROM repo
	ORDER BY name, location
`
const updateRepoHashSQL = `
	UPDATE repo SET hash=$1, modified=current_timestamp WHERE id=$2
`
const getRepoSQL = `
	SELECT id, name, location, hash, created, modified
	FROM repo
	WHERE id=$1
`

// RepoCreate creates a new repository to watch
func (db *DB) RepoCreate(name, repo string) (string, error) {
	id := xid.New()
	_, err := db.Exec(addRepoSQL, id.String(), name, repo)
	return id.String(), err
}

// RepoList get the list of repos
func (db *DB) RepoList() ([]domain.Repo, error) {
	records := []domain.Repo{}
	rows, err := db.Query(listRepoSQL)
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		r := domain.Repo{}
		err := rows.Scan(&r.ID, &r.Name, &r.Repo, &r.LastCommit, &r.Created, &r.Modified)
		if err != nil {
			return records, err
		}
		records = append(records, r)
	}

	return records, nil
}

// RepoUpdateHash updates a repo's last commit hash
func (db *DB) RepoUpdateHash(id, hash string) error {
	_, err := db.Exec(updateRepoHashSQL, hash, id)
	return err
}

// RepoGet fetches a repo from its ID
func (db *DB) RepoGet(id string) (domain.Repo, error) {
	r := domain.Repo{}
	err := db.QueryRow(getRepoSQL, id).Scan(&r.ID, &r.Name, &r.Repo, &r.LastCommit, &r.Created, &r.Modified)
	switch {
	case err == sql.ErrNoRows:
		return r, err
	case err != nil:
		log.Printf("Error retrieving database repo: %v\n", err)
		return r, err
	}
	return r, nil
}
