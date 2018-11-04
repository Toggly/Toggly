package api

import (
	"errors"

	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
)

var (
	// ErrObjectNotFound error
	ErrObjectNotFound = errors.New("object not found")
	// ErrObjectHasInheritors error
	ErrObjectHasInheritors = errors.New("object has inheritors")
)

// ObjectAPI servers object api namespace
type ObjectAPI struct {
	Owner          string
	ProjectCode    domain.ProjectCode
	EnvCode        domain.EnvironmentCode
	Storage        *storage.DataStorage
	EnvironmentAPI *EnvironmentAPI
}

func (o *ObjectAPI) envExists() error {
	_, err := o.EnvironmentAPI.Get(o.EnvCode)
	return err
}

func (o *ObjectAPI) storage() storage.ObjectStorage {
	return (*o.Storage).ForOwner(o.Owner).Projects().For(o.ProjectCode).Environments().For(o.EnvCode).Objects()
}

func (o *ObjectAPI) getInherits(obj *domain.Object) (*domain.Object, error) {
	return obj, nil
}

//List returns list of objects
func (o *ObjectAPI) List() (objects []*domain.Object, err error) {
	if err := o.envExists(); err != nil {
		return nil, err
	}
	objList, err := o.storage().List()
	if err != nil {
		return nil, err
	}
	computedObjects := make([]*domain.Object, len(objList))
	for i, obj := range objList {
		o1, e := o.getInherits(obj)
		if e != nil {
			return nil, e
		}
		computedObjects[i] = o1
	}
	return computedObjects, err
}

// Get returns object by code
func (o *ObjectAPI) Get(code domain.ObjectCode) (*domain.Object, error) {
	if err := o.envExists(); err != nil {
		return nil, err
	}
	obj, err := o.storage().Get(code)
	if err == storage.ErrNotFound {
		return nil, ErrObjectNotFound
	}
	computedObject, err := o.getInherits(obj)
	if err != nil {
		return nil, err
	}
	return computedObject, nil
}

// Create object
func (o *ObjectAPI) Create(obj *domain.Object) (*domain.Object, error) {
	if err := o.envExists(); err != nil {
		return nil, err
	}
	newObj := &domain.Object{
		Code:        obj.Code,
		Description: obj.Description,
		Inherits:    obj.Inherits,
		Parameters:  obj.Parameters,
	}
	if err := o.storage().Save(newObj); err != nil {
		return nil, err
	}
	return newObj, nil
}

// Update object
func (o *ObjectAPI) Update(obj *domain.Object) (*domain.Object, error) {
	if err := o.envExists(); err != nil {
		return nil, err
	}
	newObj := &domain.Object{
		Code:        obj.Code,
		Description: obj.Description,
		Inherits:    obj.Inherits,
		Parameters:  obj.Parameters,
	}
	if err := o.storage().Update(newObj); err != nil {
		return nil, err
	}
	return newObj, nil
}

// Delete object
func (o *ObjectAPI) Delete(code domain.ObjectCode) error {
	if err := o.envExists(); err != nil {
		return err
	}
	// TODO: check if inheritors exist
	err := o.storage().Delete(code)
	if err == storage.ErrNotFound {
		return ErrObjectNotFound
	}
	return err
}
