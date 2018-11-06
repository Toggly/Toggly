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
	// ErrObjectParentNotExists error
	ErrObjectParentNotExists = errors.New("object parrent does not exists")
	// ErrObjectInheritorTypeMismatch error
	ErrObjectInheritorTypeMismatch = errors.New("object inheritor type mismatch")
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

func (o *ObjectAPI) getParentObject(parent *domain.ObjectInheritance) (*domain.Object, error) {
	_, err := (*o.Storage).ForOwner(o.Owner).Projects().Get(parent.ProjectCode)
	if err == storage.ErrNotFound {
		return nil, ErrProjectNotFound
	}
	_, err = (*o.Storage).ForOwner(o.Owner).Projects().For(parent.ProjectCode).Environments().Get(parent.EnvCode)
	if err == storage.ErrNotFound {
		return nil, ErrEnvironmentNotFound
	}
	obj, err := (*o.Storage).ForOwner(o.Owner).Projects().For(parent.ProjectCode).Environments().For(parent.EnvCode).Objects().Get(parent.ObjectCode)
	if err == storage.ErrNotFound {
		return nil, ErrObjectNotFound
	}
	return obj, err
}

func (o *ObjectAPI) getInherits(obj *domain.Object) (*domain.Object, error) {
	if obj.Inherits == nil {
		return obj, nil
	}
	iObj, err := o.getParentObject(obj.Inherits)
	if err != nil {
		return nil, err
	}
	if iObj.Inherits != nil {
		iObj, err = o.getInherits(iObj)
		if err != nil {
			return nil, err
		}
	}
	resultParams := make([]*domain.Parameter, 0)

	for _, p := range obj.Parameters {
		for _, ip := range iObj.Parameters {
			if ip.Code == p.Code {
				if ip.Type != p.Type {
					return nil, ErrObjectInheritorTypeMismatch
				}
				resultParams = append(resultParams, &domain.Parameter{
					Code:        ip.Code,
					Description: ip.Description,
					Type:        ip.Type,
					Value:       p.Value,
				})
			}
		}
	}

	resultParams = mergeParameters(resultParams, iObj.Parameters)
	resultParams = mergeParameters(resultParams, obj.Parameters)

	return &domain.Object{
		Code:        obj.Code,
		Description: obj.Description,
		ProjectCode: obj.ProjectCode,
		EnvCode:     obj.EnvCode,
		Inherits:    obj.Inherits,
		Owner:       obj.Owner,
		Parameters:  resultParams,
	}, nil
}

func mergeParameters(arr1 []*domain.Parameter, arr2 []*domain.Parameter) []*domain.Parameter {
	res := make([]*domain.Parameter, len(arr1))
	copy(res, arr1)
	for _, p := range arr2 {
		exists := false
		for _, rp := range res {
			if rp.Code == p.Code {
				exists = true
				break
			}
		}
		if !exists {
			res = append(res, p)
		}
	}
	return res
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
		obj, e := o.getInherits(obj)
		if e != nil {
			return nil, e
		}
		computedObjects[i] = obj
	}
	return computedObjects, err
}

// Get returns object by code
func (o *ObjectAPI) Get(code domain.ObjectCode) (obj *domain.Object, err error) {
	if err = o.envExists(); err != nil {
		return nil, err
	}
	obj, err = o.storage().Get(code)
	if err == storage.ErrNotFound {
		return nil, ErrObjectNotFound
	}
	obj, err = o.getInherits(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// Create object
func (o *ObjectAPI) Create(code domain.ObjectCode, description string, inherits *domain.ObjectInheritance, parameters []*domain.Parameter) (*domain.Object, error) {
	if err := o.envExists(); err != nil {
		return nil, err
	}
	var parent *domain.Object
	var err error
	if parent, err = o.checkInheritance(inherits); err != nil {
		return nil, err
	}
	if err := o.checkParametersInheritance(parent, parameters); err != nil {
		return nil, err
	}
	newObj := &domain.Object{
		Code:        code,
		Owner:       o.Owner,
		EnvCode:     o.EnvCode,
		ProjectCode: o.ProjectCode,
		Description: description,
		Inherits:    inherits,
		Parameters:  parameters,
	}
	if err := o.storage().Save(newObj); err != nil {
		return nil, err
	}
	return o.Get(code)
}

// Update object
func (o *ObjectAPI) Update(code domain.ObjectCode, description string, inherits *domain.ObjectInheritance, parameters []*domain.Parameter) (*domain.Object, error) {
	if err := o.envExists(); err != nil {
		return nil, err
	}
	_, err := o.Get(code)
	if err != nil {
		return nil, err
	}
	var parent *domain.Object
	if parent, err = o.checkInheritance(inherits); err != nil {
		return nil, err
	}
	if err := o.checkParametersInheritance(parent, parameters); err != nil {
		return nil, err
	}
	newObj := &domain.Object{
		Code:        code,
		Owner:       o.Owner,
		EnvCode:     o.EnvCode,
		ProjectCode: o.ProjectCode,
		Description: description,
		Inherits:    inherits,
		Parameters:  parameters,
	}
	if err := o.storage().Update(newObj); err != nil {
		return nil, err
	}
	return o.Get(code)
}

func (o *ObjectAPI) checkInheritance(inherits *domain.ObjectInheritance) (*domain.Object, error) {
	if inherits == nil {
		return nil, nil
	}
	obj, err := o.getParentObject(inherits)
	if err != nil {
		switch err {
		case ErrProjectNotFound:
			return nil, ErrObjectParentNotExists
		case ErrEnvironmentNotFound:
			return nil, ErrObjectParentNotExists
		case ErrObjectNotFound:
			return nil, ErrObjectParentNotExists
		default:
			return nil, err
		}
	}
	return obj, nil
}

func (o *ObjectAPI) checkParametersInheritance(parent *domain.Object, parameters []*domain.Parameter) error {
	if parent == nil || len(parameters) == 0 {
		return nil
	}
	for _, p := range parameters {
		for _, pp := range parent.Parameters {
			if p.Code == pp.Code && p.Type != pp.Type {
				return ErrObjectInheritorTypeMismatch
			}
		}
	}
	return nil
}

// Delete object
func (o *ObjectAPI) Delete(code domain.ObjectCode) (err error) {
	if err := o.envExists(); err != nil {
		return err
	}
	inheritors, err := o.storage().ListInheritors(code)
	if err != nil {
		return err
	}
	if len(inheritors) > 0 {
		return ErrObjectHasInheritors
	}
	err = o.storage().Delete(code)
	if err == storage.ErrNotFound {
		return ErrObjectNotFound
	}
	return err
}
