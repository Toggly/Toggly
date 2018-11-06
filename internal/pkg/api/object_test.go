package api_test

import (
	"testing"

	"github.com/Toggly/core/internal/pkg/storage"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/pkg/api"
	asserts "github.com/stretchr/testify/assert"
)

const objCode = domain.ObjectCode("obj_1")

func TestObjectsWithNoProject(t *testing.T) {
	assert := asserts.New(t)

	BeforeTest()

	pApi := GetApi()
	objApi := pApi.For(ProjectCode).Environments().For(envCode).Objects()

	objList, err := objApi.List()
	assert.Equal(api.ErrProjectNotFound, err)
	assert.Nil(objList)

	env, err := objApi.Get(objCode)
	assert.Equal(api.ErrProjectNotFound, err)
	assert.Nil(env)

	env, err = objApi.Create(objCode, "", nil, nil)
	assert.Equal(api.ErrProjectNotFound, err)
	assert.Nil(env)

	env, err = objApi.Update(objCode, "", nil, nil)
	assert.Equal(api.ErrProjectNotFound, err)
	assert.Nil(env)

	err = objApi.Delete(objCode)
	assert.Equal(api.ErrProjectNotFound, err)

	AfterTest()
}

func TestObjectsWithNoEnvironment(t *testing.T) {
	assert := asserts.New(t)

	BeforeTest()

	pApi := GetApi()
	objApi := pApi.For(ProjectCode).Environments().For(envCode).Objects()

	pApi.Create(ProjectCode, "", domain.ProjectStatusActive)

	objList, err := objApi.List()
	assert.Equal(api.ErrEnvironmentNotFound, err)
	assert.Nil(objList)

	env, err := objApi.Get(objCode)
	assert.Equal(api.ErrEnvironmentNotFound, err)
	assert.Nil(env)

	env, err = objApi.Create(objCode, "", nil, nil)
	assert.Equal(api.ErrEnvironmentNotFound, err)
	assert.Nil(env)

	env, err = objApi.Update(objCode, "", nil, nil)
	assert.Equal(api.ErrEnvironmentNotFound, err)
	assert.Nil(env)

	err = objApi.Delete(objCode)
	assert.Equal(api.ErrEnvironmentNotFound, err)

	AfterTest()
}

func TestObjectsCreateErrors(t *testing.T) {
	assert := asserts.New(t)
	var err error

	BeforeTest()

	pApi := GetApi()
	envApi := pApi.For(ProjectCode).Environments()
	objApi := envApi.For(envCode).Objects()

	pApi.Create(ProjectCode, "", domain.ProjectStatusActive)
	envApi.Create(envCode, "", false)

	objList, err := objApi.List()
	assert.Nil(err)
	assert.Empty(objList)

	obj, err := objApi.Get("obj1")
	assert.Equal(api.ErrObjectNotFound, err)
	assert.Nil(obj)

	obj, err = objApi.Update("obj1", "", nil, nil)
	assert.Equal(api.ErrObjectNotFound, err)
	assert.Nil(obj)

	err = objApi.Delete("obj1")
	assert.Equal(api.ErrObjectNotFound, err)

	objApi.Create("obj1", "Obj description", nil, nil)

	_, err = objApi.Create("obj2", "", &domain.ObjectInheritance{ProjectCode: "proj2", EnvCode: envCode, ObjectCode: "obj1"}, nil)
	assert.Equal(api.ErrObjectParentNotExists, err)

	_, err = objApi.Create("obj2", "", &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: "proj123", ObjectCode: "obj1"}, nil)
	assert.Equal(api.ErrObjectParentNotExists, err)

	obj, err = objApi.Create("obj2", "", &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: envCode, ObjectCode: "obj123"}, nil)
	assert.Equal(api.ErrObjectParentNotExists, err)
	assert.Nil(obj)

	AfterTest()
}

func TestObjects(t *testing.T) {
	assert := asserts.New(t)

	BeforeTest()

	pApi := GetApi()
	envApi := pApi.For(ProjectCode).Environments()
	objApi := envApi.For(envCode).Objects()

	pApi.Create(ProjectCode, "", domain.ProjectStatusActive)
	envApi.Create(envCode, "", false)

	obj, err := objApi.Create("obj1", "Obj description", nil, nil)
	assert.Nil(err)
	assert.NotNil(obj)
	assert.Equal(ow, obj.Owner)
	assert.Equal(ProjectCode, obj.ProjectCode)
	assert.Equal(envCode, obj.EnvCode)
	assert.Equal(domain.ObjectCode("obj1"), obj.Code)
	assert.Equal("Obj description", obj.Description)
	assert.Nil(obj.Inherits)
	assert.Empty(obj.Parameters)

	_, err = objApi.Create("obj1", "Obj description", nil, nil)
	assert.IsType(&storage.UniqueIndexError{}, err)

	objList, err := objApi.List()
	assert.Nil(err)
	assert.Len(objList, 1)

	obj, err = objApi.Get("obj1")
	assert.Nil(err)
	assert.NotNil(obj)
	assert.Equal(ow, obj.Owner)
	assert.Equal(ProjectCode, obj.ProjectCode)
	assert.Equal(envCode, obj.EnvCode)
	assert.Equal(domain.ObjectCode("obj1"), obj.Code)
	assert.Equal("Obj description", obj.Description)
	assert.Nil(obj.Inherits)
	assert.Empty(obj.Parameters)

	params := []*domain.Parameter{
		&domain.Parameter{
			Code:        "param1",
			Description: "Param 1",
			Type:        domain.ParameterBool,
			Value:       true,
		},
	}

	obj, err = objApi.Update("obj1", "New description", nil, params)
	assert.Nil(err)
	assert.NotNil(obj)
	assert.Equal(ow, obj.Owner)
	assert.Equal(ProjectCode, obj.ProjectCode)
	assert.Equal(envCode, obj.EnvCode)
	assert.Equal(domain.ObjectCode("obj1"), obj.Code)
	assert.Equal("New description", obj.Description)
	assert.Nil(obj.Inherits)
	assert.Equal(params, obj.Parameters)

	AfterTest()
}

func TestObjectsDeleteWithInheritors(t *testing.T) {
	assert := asserts.New(t)
	var err error

	BeforeTest()

	pApi := GetApi()
	envApi := pApi.For(ProjectCode).Environments()
	objApi := envApi.For(envCode).Objects()

	pApi.Create(ProjectCode, "", domain.ProjectStatusActive)
	envApi.Create(envCode, "", false)

	objApi.Create("obj1", "", nil, nil)

	_, err = objApi.Create("obj2", "", &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: envCode, ObjectCode: "obj1"}, nil)
	assert.Nil(err)

	err = objApi.Delete("obj1")
	assert.Equal(api.ErrObjectHasInheritors, err)

	err = objApi.Delete("obj2")
	assert.Nil(err)

	err = objApi.Delete("obj1")
	assert.Nil(err)

	AfterTest()
}

func TestObjectsInheritance(t *testing.T) {
	assert := asserts.New(t)

	BeforeTest()

	pApi := GetApi()
	envApi := pApi.For(ProjectCode).Environments()
	objApi := envApi.For(envCode).Objects()

	pApi.Create(ProjectCode, "", domain.ProjectStatusActive)
	envApi.Create(envCode, "", false)

	parentParams := []*domain.Parameter{
		&domain.Parameter{
			Code:        "param1",
			Description: "Param 1",
			Type:        domain.ParameterBool,
			Value:       true,
		},
		&domain.Parameter{
			Code:        "param2",
			Description: "Param 2",
			Type:        domain.ParameterBool,
			Value:       true,
		},
	}

	childParams := []*domain.Parameter{
		&domain.Parameter{
			Code:        "param2",
			Description: "Param 2 Child has to be overridden by original",
			Type:        domain.ParameterBool,
			Value:       false,
		},
		&domain.Parameter{
			Code:        "param3",
			Description: "Param 3",
			Type:        domain.ParameterString,
			Value:       "value",
		},
	}

	objApi.Create("obj1", "", nil, parentParams)
	objApi.Create("obj2", "", &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: envCode, ObjectCode: "obj1"}, childParams)

	obj, err := objApi.Get("obj2")
	assert.Nil(err)
	assert.NotNil(obj)
	assert.Len(obj.Parameters, 3)

	for _, p := range obj.Parameters {
		switch p.Code {
		case "param1":
			assert.Equal("Param 1", p.Description)
			assert.Equal(domain.ParameterBool, p.Type)
			assert.Equal(true, p.Value)
		case "param2":
			assert.Equal("Param 2", p.Description)
			assert.Equal(domain.ParameterBool, p.Type)
			assert.Equal(false, p.Value)
		case "param3":
			assert.Equal("Param 3", p.Description)
			assert.Equal(domain.ParameterString, p.Type)
			assert.Equal("value", p.Value)
		}
	}

	wrongParams := []*domain.Parameter{
		&domain.Parameter{
			Code:        "param2",
			Description: "Param 2",
			Type:        domain.ParameterString,
			Value:       "wrong type",
		},
	}
	_, err = objApi.Create("obj3", "", &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: envCode, ObjectCode: "obj1"}, wrongParams)
	assert.Equal(api.ErrObjectInheritorTypeMismatch, err)

	obj3Params := []*domain.Parameter{
		&domain.Parameter{
			Code:        "param3",
			Description: "Param 3 will be overridden",
			Type:        domain.ParameterString,
			Value:       "value 2",
		},
	}

	obj, err = objApi.Create("obj3", "", &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: envCode, ObjectCode: "obj2"}, obj3Params)
	assert.Nil(err)
	assert.NotNil(obj)
	assert.Len(obj.Parameters, 3)

	for _, p := range obj.Parameters {
		switch p.Code {
		case "param1":
			assert.Equal("Param 1", p.Description)
			assert.Equal(domain.ParameterBool, p.Type)
			assert.Equal(true, p.Value)
		case "param2":
			assert.Equal("Param 2", p.Description)
			assert.Equal(domain.ParameterBool, p.Type)
			assert.Equal(false, p.Value)
		case "param3":
			assert.Equal("Param 3", p.Description)
			assert.Equal(domain.ParameterString, p.Type)
			assert.Equal("value 2", p.Value)
		}
	}

	AfterTest()

}
