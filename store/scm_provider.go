package store

import (
	"github.com/dhruv1397/prm/types"
)

type SCMProvider interface {
	Create(provider types.SCMProvider) error
	UpdateBulk(providers []types.SCMProvider) error
	List(providerType string, providerName string) ([]*types.SCMProvider, error)
	Delete(name string) error
	Purge() error
}
