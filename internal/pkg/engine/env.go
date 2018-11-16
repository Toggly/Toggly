package engine

import (
	"time"

	"github.com/Toggly/core/internal/api"
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

func (e *EnvironmentAPI) projectExists() error {
	_, err := e.ProjectAPI.Get(e.ProjectCode)
	return err
}

func (e *EnvironmentAPI) storage() storage.EnvironmentStorage {
	return (*e.Storage).ForOwner(e.Owner).Projects().For(e.ProjectCode).Environments()
}

// List returns list of project environments
func (e *EnvironmentAPI) List() ([]*domain.Environment, error) {
	if err := e.projectExists(); err != nil {
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
	if err := e.projectExists(); err != nil {
		return nil, err
	}
	env, err := e.storage().Get(code)
	if err == storage.ErrNotFound {
		return nil, api.ErrEnvironmentNotFound
	}
	return env, err
}

func checkEnvParams(code domain.EnvironmentCode, description string, protected bool) error {
	if code == "" {
		return api.NewBadRequestError("Environment code not specified")
	}
	return nil
}

// Create environment
func (e *EnvironmentAPI) Create(info *api.EnvironmentInfo) (*domain.Environment, error) {
	code := info.Code
	description := info.Description
	protected := info.Protected
	if err := e.projectExists(); err != nil {
		return nil, err
	}
	if err := checkEnvParams(code, description, protected); err != nil {
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

// Update environment
func (e *EnvironmentAPI) Update(info *api.EnvironmentInfo) (*domain.Environment, error) {
	code := info.Code
	description := info.Description
	protected := info.Protected
	if err := e.projectExists(); err != nil {
		return nil, err
	}
	if err := checkEnvParams(code, description, protected); err != nil {
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

// Delete environment
func (e *EnvironmentAPI) Delete(code domain.EnvironmentCode) error {
	if err := e.projectExists(); err != nil {
		return err
	}
	objList, err := e.For(code).Objects().List()
	if err != nil {
		return err
	}
	if len(objList) > 0 {
		return api.ErrEnvironmentNotEmpty
	}
	err = e.storage().Delete(code)
	if err == storage.ErrNotFound {
		return api.ErrEnvironmentNotFound
	}
	return err
}

// For returns object api for specified environment
func (e *EnvironmentAPI) For(code domain.EnvironmentCode) api.ForObjectAPI {
	return &ForObjectAPI{
		Owner:          e.Owner,
		ProjectCode:    e.ProjectCode,
		EnvCode:        code,
		Storage:        e.Storage,
		EnvironmentAPI: e,
	}
}

// ForObjectAPI type
type ForObjectAPI struct {
	Owner          string
	ProjectCode    domain.ProjectCode
	EnvCode        domain.EnvironmentCode
	Storage        *storage.DataStorage
	EnvironmentAPI *EnvironmentAPI
}

// Objects returns ObjectAPI
func (fo *ForObjectAPI) Objects() api.ObjectAPI {
	return &ObjectAPI{
		Owner:          fo.Owner,
		ProjectCode:    fo.ProjectCode,
		EnvCode:        fo.EnvCode,
		Storage:        fo.Storage,
		EnvironmentAPI: fo.EnvironmentAPI,
	}
}
