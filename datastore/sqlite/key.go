package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/ogra1/fabrica/domain"
	"github.com/rs/xid"
	"log"
)

const createKeysTableSQL = `
	CREATE TABLE IF NOT EXISTS keys (
		id               varchar(200) primary key not null,
		name             varchar(200) UNIQUE not null,
		username         varchar(200) not null,
		data             text not null,
		password         varchar(200) default '',
		created          timestamp default current_timestamp,
		modified         timestamp default current_timestamp
	)
`
const addKeysSQL = `
	INSERT INTO keys (id, name, username, data, password) VALUES ($1, $2, $3, $4, $5)
`
const getKeysSQL = `
	SELECT id, name, username, data, password
	FROM keys
	WHERE id=$1
`
const listKeysSQL = `
	SELECT id, name, username, created
	FROM keys
	ORDER BY name
`
const deleteKeysSQL = `
	DELETE FROM keys WHERE id=$1
`

// KeysCreate stores a new ssh key
func (db *DB) KeysCreate(name, username, data, password string) (string, error) {
	// Get the secret key
	secret, err := db.secretKey()
	if err != nil {
		return "", err
	}

	// Encrypt the secret key
	dataEnc, err := encryptKey(secret, data)
	if err != nil {
		return "", err
	}
	passwordEnc, err := encryptKey(secret, password)
	if err != nil {
		return "", err
	}

	log.Println("---", name, username, data, password)

	// Save the encrypted record
	id := xid.New()
	log.Println("---", id.String(), name, username, dataEnc, passwordEnc)
	_, err = db.Exec(addKeysSQL, id.String(), name, username, dataEnc, passwordEnc)
	return id.String(), err
}

// KeysGet fetches an ssh key by its ID
func (db *DB) KeysGet(id string) (domain.Key, error) {
	r := domain.Key{}
	err := db.QueryRow(getKeysSQL, id).Scan(&r.ID, &r.Name, &r.Username, &r.Data, &r.Password)
	switch {
	case err == sql.ErrNoRows:
		return r, err
	case err != nil:
		log.Printf("Error retrieving database repo: %v\n", err)
		return r, err
	}

	// Get the secret key
	secret, err := db.secretKey()
	if err != nil {
		return r, err
	}

	// Decrypt the data
	r.Data, err = decryptKey([]byte(r.Data), secret)
	if err != nil {
		return r, err
	}
	r.Password, err = decryptKey([]byte(r.Password), secret)
	if err != nil {
		return r, err
	}

	return r, nil
}

// KeysList get the list of ssh keys. Only unencrypted data is returned.
func (db *DB) KeysList() ([]domain.Key, error) {
	records := []domain.Key{}
	rows, err := db.Query(listKeysSQL)
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		r := domain.Key{}
		err := rows.Scan(&r.ID, &r.Name, &r.Username, &r.Created)
		if err != nil {
			return records, err
		}
		records = append(records, r)
	}

	return records, nil
}

// KeysDelete removes a key from its name
func (db *DB) KeysDelete(id string) error {
	// Check the key is not used
	repos, err := db.ReposForKey(id)
	if err != nil {
		return err
	}
	if len(repos) > 0 {
		return fmt.Errorf("the key is used by %d repositories", len(repos))
	}

	_, err = db.Exec(deleteKeysSQL, id)
	return err
}
