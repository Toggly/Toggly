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

func (e *EnvironmentAPI) storage() storage.EnvironmentStorage {
	return (*e.Storage).ForOwner(e.Owner).Projects().For(e.ProjectCode).Environments()
}

// List returns list of project environments
func (e *EnvironmentAPI) List() ([]*domain.Environment, error) {
	if err := e.projectExists(e.ProjectCode); err != nil {
		return nil, err
	}
	envList, err := e.storage().List()
	if err != nil {
		return nil, err
	}
	return envList, err
}

// Get returns environment by code
func (e *EnvironmentAPI) Get(code domain.EnvironmentCode) (*domain.Environment, error) {
	if err := e.projectExists(e.ProjectCode); err != nil {
		return nil, err
	}
	env, err := e.storage().Get(code)
	if err == storage.ErrNotFound {
		return nil, ErrEnvironmentNotFound
	}
	return env, err
}

// Create saves environment
func (e *EnvironmentAPI) Create(code domain.EnvironmentCode, description string, protected bool) (*domain.Environment, error) {
	if err := e.projectExists(e.ProjectCode); err != nil {
		return nil, err
	}
	newEnv := &domain.Environment{
		Code:        code,
		Description: description,
		OwnerID:     e.Owner,
		ProjectCode: e.ProjectCode,
		Protected:   protected,
		RegDate:     bson.Now().In(time.UTC),
	}
	if err := e.storage().Save(newEnv); err != nil {
		return nil, err
	}
	return newEnv, nil
}

// Update saves environment
func (e *EnvironmentAPI) Update(code domain.EnvironmentCode, description string, protected bool) (*domain.Environment, error) {
	if err := e.projectExists(e.ProjectCode); err != nil {
		return nil, err
	}
	uEnv, err := e.Get(code)
	if err != nil {
		return nil, err
	}
	newEnv := &domain.Environment{
		Code:        code,
		Description: description,
		OwnerID:     e.Owner,
		ProjectCode: e.ProjectCode,
		Protected:   protected,
		RegDate:     uEnv.RegDate,
	}
	if err := e.storage().Update(newEnv); err != nil {
		return nil, err
	}
	return newEnv, nil
}

// Delete Environment
func (e *EnvironmentAPI) Delete(code domain.EnvironmentCode) error {
	if err := e.projectExists(e.ProjectCode); err != nil {
		return err
	}
	// TODO: check if env is empty
	err := e.storage().Delete(code)
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
