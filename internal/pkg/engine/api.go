package engine

import (
	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/pkg/storage"
)

// NewTogglyAPI returns api engine
func NewTogglyAPI(storage *storage.DataStorage) api.TogglyAPI {
	return &Engine{Storage: storage}
}

// Engine type
type Engine struct {
	Storage *storage.DataStorage
}

// ForOwner returns owner api
func (e *Engine) ForOwner(owner string) api.OwnerAPI {
	return &OwnerAPI{Owner: owner, Storage: e.Storage}
}

// OwnerAPI type
type OwnerAPI struct {
	Owner   string
	Storage *storage.DataStorage
}

// Projects returns project api
func (o *OwnerAPI) Projects() api.ProjectAPI {
	return &ProjectAPI{*o}
}
