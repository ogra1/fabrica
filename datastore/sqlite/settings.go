package sqlite

import (
	"database/sql"
	"github.com/ogra1/fabrica/domain"
	"github.com/rs/xid"
	"log"
)

const (
	secretKeyName = "secret"
	keyLength     = 64
)

const createSettingsTableSQL string = `
	CREATE TABLE IF NOT EXISTS settings (
		id               varchar(200) primary key not null,
		key              varchar(200) not null,
		name             varchar(200) not null,
		data             text default '',
		created          timestamp default current_timestamp,
		modified         timestamp default current_timestamp
	)
`
const addSettingSQL = `
	INSERT INTO settings (id, key, name, data) VALUES ($1, $2, $3, $4)
`
const getSettingSQL = `
	SELECT id, key, name, data
	FROM settings
	WHERE key=$1 and name=$2
`

// SettingsCreate stores a new config setting
func (db *DB) SettingsCreate(key, name, data string) (string, error) {
	id := xid.New()
	_, err := db.Exec(addSettingSQL, id.String(), key, name, data)
	return id.String(), err
}

// SettingsGet fetches an existing config setting
func (db *DB) SettingsGet(key, name string) (domain.ConfigSetting, error) {
	r := domain.ConfigSetting{}
	err := db.QueryRow(getSettingSQL, key, name).Scan(&r.ID, &r.Key, &r.Name, &r.Data)
	switch {
	case err == sql.ErrNoRows:
		return r, err
	case err != nil:
		log.Printf("Error retrieving database repo: %v\n", err)
		return r, err
	}
	return r, nil
}

func (db *DB) secretKey() (string, error) {
	// Get the secret key from the database and return it
	key, err := db.SettingsGet(secretKeyName, secretKeyName)
	if err == nil {
		return key.Data, nil
	}

	// Cannot find a secret, so generate one
	secret, err := generateSecret()
	if err != nil {
		return "", err
	}

	// Store the secret
	if _, err := db.SettingsCreate(secretKeyName, secretKeyName, secret); err != nil {
		return "", err
	}
	return secret, nil
}
