package key

import (
	"github.com/ogra1/fabrica/datastore"
	"github.com/ogra1/fabrica/domain"
)

// Srv is the interface for the ssh key service
type Srv interface {
	Create(name, username, data, password string) (string, error)
	Get(id string) (domain.Key, error)
	List() ([]domain.Key, error)
	Delete(id string) error
}

// Service implements the ssh key service
type Service struct {
	Datastore datastore.Datastore
}

// NewKeyService creates a new ssh key service
func NewKeyService(ds datastore.Datastore) *Service {
	return &Service{
		Datastore: ds,
	}
}

// Create stores a new ssh key
func (ks *Service) Create(name, username, data, password string) (string, error) {
	return ks.Datastore.KeysCreate(name, username, data, password)
}

// Get fetches an existing ssh key
func (ks *Service) Get(id string) (domain.Key, error) {
	return ks.Datastore.KeysGet(id)
}

// List fetches existing ssh keys
func (ks *Service) List() ([]domain.Key, error) {
	return ks.Datastore.KeysList()
}

// Delete removes an existing ssh key
func (ks *Service) Delete(id string) error {
	return ks.Datastore.KeysDelete(id)
}
