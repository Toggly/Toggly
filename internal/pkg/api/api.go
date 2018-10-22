package api

import (
	"errors"

	"github.com/Toggly/core/internal/pkg/storage"
)

var (
	// ErrNotFound error
	ErrNotFound = errors.New("not found")
)

// Engine type
type Engine struct {
	Storage *storage.DataStorage
}

// ForOwner returns owner api
func (e *Engine) ForOwner(owner string) *OwnerAPI {
	return &OwnerAPI{Owner: owner, Storage: e.Storage}
}

// OwnerAPI type
type OwnerAPI struct {
	Owner   string
	Storage *storage.DataStorage
}

// Project returns project api
func (o *OwnerAPI) Project() *ProjectAPI {
	return &ProjectAPI{Owner: o.Owner, Storage: o.Storage}
}
