package sqlite

import (
	"github.com/ogra1/fabrica/domain"
	"github.com/rs/xid"
)

const createBuildLogTableSQL string = `
	CREATE TABLE IF NOT EXISTS buildlog (
		id               varchar(200) primary key not null,
		build_id         varchar(200) not null,
		message          text,
		created          timestamp default current_timestamp,
		FOREIGN KEY (build_id) REFERENCES build (id)
	)
`
const addBuildLogSQL = `
	INSERT INTO buildlog(id, build_id, message) VALUES ($1, $2, $3)
`
const listBuildLogSQL = `
	SELECT id, build_id, message, created
	FROM buildlog
	WHERE build_id=$1
	ORDER BY created
`
const deleteBuildLogsSQL = `
	DELETE FROM buildlog WHERE build_id=$1
`

// BuildLogCreate logs a message for a build
func (db *DB) BuildLogCreate(buildID, message string) error {
	id := xid.New()
	_, err := db.Exec(addBuildLogSQL, id.String(), buildID, message)
	return err
}

// BuildLogList lists messages for a build
func (db *DB) BuildLogList(buildID string) ([]domain.BuildLog, error) {
	logs := []domain.BuildLog{}
	rows, err := db.Query(listBuildLogSQL, buildID)
	if err != nil {
		return logs, err
	}
	defer rows.Close()

	for rows.Next() {
		r := domain.BuildLog{}
		err := rows.Scan(&r.ID, &r.BuildID, &r.Message, &r.Created)
		if err != nil {
			return logs, err
		}
		logs = append(logs, r)
	}

	return logs, nil
}

// BuildLogDelete deletes logs for a build
func (db *DB) BuildLogDelete(buildID string) error {
	_, err := db.Exec(deleteBuildLogsSQL, buildID)
	return err
}
