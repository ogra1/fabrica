package service

import (
	"encoding/base64"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/ogra1/fabrica/domain"
	ssh2 "golang.org/x/crypto/ssh"
	"log"
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

// GitAuth returns the ssh auth method for git
func GitAuth(key domain.Key, gitURL string) (transport.AuthMethod, error) {
	// Decode the private key
	var data []byte
	data, err := base64.StdEncoding.DecodeString(key.Data)
	if err != nil {
		log.Println("Error decoding ssh key:", err)
		return nil, err
	}

	// Get the ssh username from the URL
	username, err := usernameFromRepo(gitURL)
	if err != nil {
		log.Println("Error parsing git URL:", err)
		return nil, err
	}

	// Set the ssh auth for git
	pubKeys, err := ssh.NewPublicKeys(username, data, key.Password)
	if err != nil {
		log.Println("Error creating ssh key auth:", err)
		return nil, err
	}

	// Disable the known_hosts check
	pubKeys.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
	return pubKeys, nil
}

func usernameFromRepo(u string) (string, error) {
	endpoint, err := transport.NewEndpoint(u)
	if err != nil {
		return "", err
	}

	return endpoint.User, nil
}
