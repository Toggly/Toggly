package engine

import (
	"fmt"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
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
		return nil, api.ErrProjectNotFound
	}
	_, err = (*o.Storage).ForOwner(o.Owner).Projects().For(parent.ProjectCode).Environments().Get(parent.EnvCode)
	if err == storage.ErrNotFound {
		return nil, api.ErrEnvironmentNotFound
	}
	obj, err := (*o.Storage).ForOwner(o.Owner).Projects().For(parent.ProjectCode).Environments().For(parent.EnvCode).Objects().Get(parent.ObjectCode)
	if err == storage.ErrNotFound {
		return nil, api.ErrObjectNotFound
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
					return nil, api.ErrObjectInheritorTypeMismatch
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
		return nil, api.ErrObjectNotFound
	}
	obj, err = o.getInherits(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func checkObjParams(code domain.ObjectCode, description string, inherits *domain.ObjectInheritance, parameters []*domain.Parameter) error {
	if code == "" {
		return api.NewBadRequestError("Object code not specified")
	}
	if parameters != nil && len(parameters) > 0 {
		for _, p := range parameters {
			if p.Code == "" {
				return api.NewBadRequestError("Parameter code not specified")
			}
			if p.Value == nil {
				return api.NewBadRequestError("Parameter value not specified")
			}
			if !isParameterType(p.Type) {
				return api.NewBadRequestError("Parameter type not specified or wrong")
			}
			if p.AllowedValues != nil {
				err := checkAllowedValue(p.Type, p.Value, p.AllowedValues)
				if err != nil {
					return api.NewBadRequestError(err.Error())
				}
			}
		}
	}
	return nil
}

func isParameterType(t domain.ParameterType) bool {
	return t == domain.ParameterBool || t == domain.ParameterInt || t == domain.ParameterString
}

func checkAllowedValue(t domain.ParameterType, v interface{}, allowed []interface{}) error {
	switch t {
	case domain.ParameterBool:
		if _, ok := v.(bool); !ok {
			return fmt.Errorf("Can't cast parameter value `%v` to bool", v)
		}
	case domain.ParameterInt:
		if _, ok := v.(int); !ok {
			return fmt.Errorf("Can't cast parameter value `%v` to int", v)
		}
	case domain.ParameterString:
		if _, ok := v.(string); !ok {
			return fmt.Errorf("Can't cast parameter value `%v` to string", v)
		}
	}

	return nil
}

// Create object
func (o *ObjectAPI) Create(info *api.ObjectInfo) (*domain.Object, error) {
	code := info.Code
	description := info.Description
	inherits := info.Inherits
	parameters := info.Parameters
	if err := o.envExists(); err != nil {
		return nil, err
	}
	if err := checkObjParams(code, description, inherits, parameters); err != nil {
		return nil, err
	}
	var parent *domain.Object
	var err error
	if parent, err = o.checkInheritance(inherits); err != nil {
		return nil, err
	}
	if err := o.checkParametersInheritanceForParent(parent, parameters); err != nil {
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
func (o *ObjectAPI) Update(info *api.ObjectInfo) (*domain.Object, error) {
	code := info.Code
	description := info.Description
	inherits := info.Inherits
	parameters := info.Parameters
	if err := o.envExists(); err != nil {
		return nil, err
	}
	if err := checkObjParams(code, description, inherits, parameters); err != nil {
		return nil, err
	}
	obj, err := o.Get(code)
	if err != nil {
		return nil, err
	}
	if err := o.checkIfParametersChanged(obj, parameters); err != nil {
		return nil, err
	}
	var parent *domain.Object
	if parent, err = o.checkInheritance(inherits); err != nil {
		return nil, err
	}
	if err := o.checkParametersInheritanceForParent(parent, parameters); err != nil {
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

// InheritorsFlatList returns flat list of inheritors
func (o *ObjectAPI) InheritorsFlatList(code domain.ObjectCode) ([]*domain.Object, error) {
	list, err := o.storage().ListInheritors(code)
	if err != nil {
		return nil, err
	}
	for _, i := range list {
		sub, err := o.InheritorsFlatList(i.Code)
		if err != nil {
			return nil, err
		}
		list = append(list, sub...)
	}
	return list, nil
}

func (o *ObjectAPI) checkIfParametersChanged(obj *domain.Object, parameters []*domain.Parameter) error {
	inheritors, err := o.InheritorsFlatList(obj.Code)
	if err != nil {
		return err
	}
	curParams := make(map[domain.ParameterCode]*domain.Parameter)
	for _, par := range obj.Parameters {
		curParams[par.Code] = par
	}
	for _, p := range parameters {
		par := curParams[p.Code]
		if par != nil {
			// 1. Deny changing the parameter type
			if par.Type != p.Type {
				return api.NewObjectParameterError(string(par.Code), "Object parameter type changing restricted")
			}
		} else {
			// 2. Deny creating the parameter if inheritors already have it
			for _, inh := range inheritors {
				for _, ip := range inh.Parameters {
					if ip.Code == p.Code {
						msg := fmt.Sprintf("Object parameter exists in inheritor: %s:%s:%s", inh.ProjectCode, inh.EnvCode, inh.Code)
						return api.NewObjectParameterError(string(p.Code), msg)
					}
				}
			}
		}
	}
	return nil
}

func (o *ObjectAPI) checkInheritance(inherits *domain.ObjectInheritance) (*domain.Object, error) {
	if inherits == nil {
		return nil, nil
	}
	obj, err := o.getParentObject(inherits)
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			return nil, api.ErrObjectParentNotExists
		case api.ErrEnvironmentNotFound:
			return nil, api.ErrObjectParentNotExists
		case api.ErrObjectNotFound:
			return nil, api.ErrObjectParentNotExists
		default:
			return nil, err
		}
	}
	return obj, nil
}

func (o *ObjectAPI) checkParametersInheritanceForParent(parent *domain.Object, parameters []*domain.Parameter) error {
	if parent == nil || len(parameters) == 0 {
		return nil
	}
	for _, p := range parameters {
		for _, pp := range parent.Parameters {
			if p.Code == pp.Code && p.Type != pp.Type {
				return api.ErrObjectInheritorTypeMismatch
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
		return api.ErrObjectHasInheritors
	}
	err = o.storage().Delete(code)
	if err == storage.ErrNotFound {
		return api.ErrObjectNotFound
	}
	return err
}
