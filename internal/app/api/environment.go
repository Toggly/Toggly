package api

import (
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
)

// EnvironmentAPI type
type EnvironmentAPI struct {
	Owner       string
	ProjectCode domain.ProjectCode
	Storage     *storage.DataStorage
}

// List returns list of project environments
func (e *EnvironmentAPI) List() ([]*domain.Environment, error) {
	return (*e.Storage).
		ForOwner(e.Owner).
		Projects().
		For(e.ProjectCode).
		Environments().
		List()
}

// Get returns environment by code
func (e *EnvironmentAPI) Get(code domain.EnvironmentCode) (*domain.Environment, error) {
	return nil, nil
}

// Save saves environment
func (e *EnvironmentAPI) Save(env *domain.Environment) error {
	return nil
}

// For returns object api for specified environment
func (e *EnvironmentAPI) For(code domain.EnvironmentCode) *ObjectAPI {
	return &ObjectAPI{
		Owner:   e.Owner,
		Project: e.ProjectCode,
		Env:     code,
		Storage: e.Storage,
	}
}
