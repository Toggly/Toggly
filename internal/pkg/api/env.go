package api

import (
	"time"

	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/globalsign/mgo/bson"
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

// Create saves environment
func (e *EnvironmentAPI) Create(code domain.EnvironmentCode, description string, protected bool) (*domain.Environment, error) {
	newEnv := &domain.Environment{
		Code:        code,
		Description: description,
		OwnerID:     e.Owner,
		ProjectCode: e.ProjectCode,
		Protected:   protected,
		RegDate:     bson.Now().In(time.UTC),
	}
	if err := (*e.Storage).ForOwner(e.Owner).Projects().For(e.ProjectCode).Environments().Save(newEnv); err != nil {
		return nil, err
	}
	return newEnv, nil
}

// Update saves environment
func (e *EnvironmentAPI) Update(env *domain.Environment) (*domain.Environment, error) {
	uEnv, err := e.Get(env.Code)
	if err != nil {
		return nil, err
	}
	newEnv := &domain.Environment{
		Code:        env.Code,
		Description: env.Description,
		OwnerID:     e.Owner,
		ProjectCode: e.ProjectCode,
		Protected:   env.Protected,
		RegDate:     uEnv.RegDate,
	}
	if err := (*e.Storage).ForOwner(e.Owner).Projects().For(e.ProjectCode).Environments().Update(newEnv); err != nil {
		return nil, err
	}
	return env, nil
}

// Delete Environment
func (e *EnvironmentAPI) Delete(code domain.EnvironmentCode) error {
	// TODO: check if env is empty
	err := (*e.Storage).ForOwner(e.Owner).Projects().For(e.ProjectCode).Environments().Delete(code)
	if err == storage.ErrNotFound {
		return ErrEnvironmentNotFound
	}
	return err
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
