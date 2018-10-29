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
	ProjectAPI  *ProjectAPI
}

func (e *EnvironmentAPI) projectExists(p domain.ProjectCode) error {
	_, err := e.ProjectAPI.Get(e.ProjectCode)
	return err
}

// List returns list of project environments
func (e *EnvironmentAPI) List() ([]*domain.Environment, error) {

	err := e.projectExists(e.ProjectCode)
	if err != nil {
		return nil, err
	}

	envList, err := (*e.Storage).ForOwner(e.Owner).Projects().For(e.ProjectCode).Environments().List()
	if err != nil {
		return nil, err
	}

	return envList, err
}

// Get returns environment by code
func (e *EnvironmentAPI) Get(code domain.EnvironmentCode) (*domain.Environment, error) {
	return nil, nil
}

// Save saves environment
func (e *EnvironmentAPI) Save(env *domain.Environment) (*domain.Environment, error) {
	if err := (*e.Storage).ForOwner(e.Owner).Projects().For(e.ProjectCode).Environments().Save(env); err != nil {
		return nil, err
	}
	return env, nil
}

// For returns object api for specified environment
func (e *EnvironmentAPI) For(code domain.EnvironmentCode) *ObjectAPI {
	return &ObjectAPI{
		Owner:          e.Owner,
		Project:        e.ProjectCode,
		Env:            code,
		Storage:        e.Storage,
		EnvironmentAPI: e,
	}
}
