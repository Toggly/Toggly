package engine_test

import (
	"testing"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/pkg/storage"

	"github.com/Toggly/core/internal/domain"

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

	env, err = objApi.Create(&api.ObjectInfo{Code: objCode})
	assert.Equal(api.ErrProjectNotFound, err)
	assert.Nil(env)

	env, err = objApi.Update(&api.ObjectInfo{Code: objCode})
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

	pApi.Create(&api.ProjectInfo{Code: ProjectCode, Status: domain.ProjectStatusActive})

	objList, err := objApi.List()
	assert.Equal(api.ErrEnvironmentNotFound, err)
	assert.Nil(objList)

	env, err := objApi.Get(objCode)
	assert.Equal(api.ErrEnvironmentNotFound, err)
	assert.Nil(env)

	env, err = objApi.Create(&api.ObjectInfo{Code: objCode})
	assert.Equal(api.ErrEnvironmentNotFound, err)
	assert.Nil(env)

	env, err = objApi.Update(&api.ObjectInfo{Code: objCode})
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

	pApi.Create(&api.ProjectInfo{Code: ProjectCode, Status: domain.ProjectStatusActive})
	envApi.Create(&api.EnvironmentInfo{Code: envCode})

	objList, err := objApi.List()
	assert.Nil(err)
	assert.Empty(objList)

	obj, err := objApi.Get("obj1")
	assert.Equal(api.ErrObjectNotFound, err)
	assert.Nil(obj)

	obj, err = objApi.Update(&api.ObjectInfo{Code: "obj1"})
	assert.Equal(api.ErrObjectNotFound, err)
	assert.Nil(obj)

	err = objApi.Delete("obj1")
	assert.Equal(api.ErrObjectNotFound, err)

	objApi.Create(&api.ObjectInfo{
		Code:        "obj1",
		Description: "Obj description",
	})

	_, err = objApi.Create(&api.ObjectInfo{
		Code: "obj2",
		Inherits: &domain.ObjectInheritance{
			ProjectCode: "proj2",
			EnvCode:     envCode,
			ObjectCode:  "obj1",
		},
	})
	assert.Equal(api.ErrObjectParentNotExists, err)

	_, err = objApi.Create(&api.ObjectInfo{
		Code: "obj2",
		Inherits: &domain.ObjectInheritance{
			ProjectCode: ProjectCode,
			EnvCode:     "proj123",
			ObjectCode:  "obj1",
		},
	})
	assert.Equal(api.ErrObjectParentNotExists, err)

	_, err = objApi.Create(&api.ObjectInfo{
		Code: "obj2",
		Inherits: &domain.ObjectInheritance{
			ProjectCode: ProjectCode,
			EnvCode:     envCode,
			ObjectCode:  "obj123",
		},
	})
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

	pApi.Create(&api.ProjectInfo{Code: ProjectCode, Status: domain.ProjectStatusActive})
	envApi.Create(&api.EnvironmentInfo{Code: envCode})

	_, err := objApi.Create(&api.ObjectInfo{Description: "Obj description"})
	assert.IsType(&api.ErrBadRequest{}, err)

	_, err = objApi.Update(&api.ObjectInfo{Description: "Obj description"})
	assert.IsType(&api.ErrBadRequest{}, err)

	obj, err := objApi.Create(&api.ObjectInfo{
		Code:        "obj1",
		Description: "Obj description",
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:        "param1",
				Description: "Param 1",
			},
		},
	})
	assert.IsType(&api.ErrBadRequest{}, err)

	obj, err = objApi.Create(&api.ObjectInfo{
		Code:        "obj1",
		Description: "Obj description",
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:        "param1",
				Description: "Param 1",
				Type:        domain.ParameterBool,
			},
		},
	})
	assert.IsType(&api.ErrBadRequest{}, err)

	obj, err = objApi.Create(&api.ObjectInfo{
		Code:        "obj1",
		Description: "Obj description",
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:        "param1",
				Description: "Param 1",
				Type:        "wrong_type",
				Value:       "some_value",
			},
		},
	})
	assert.IsType(&api.ErrBadRequest{}, err)

	obj, err = objApi.Create(&api.ObjectInfo{Code: "obj1", Description: "Obj description"})
	assert.Nil(err)
	assert.NotNil(obj)
	assert.Equal(ow, obj.Owner)
	assert.Equal(ProjectCode, obj.ProjectCode)
	assert.Equal(envCode, obj.EnvCode)
	assert.Equal(domain.ObjectCode("obj1"), obj.Code)
	assert.Equal("Obj description", obj.Description)
	assert.Nil(obj.Inherits)
	assert.Empty(obj.Parameters)

	_, err = objApi.Create(&api.ObjectInfo{Code: "obj1"})
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

	obj, err = objApi.Update(&api.ObjectInfo{
		Code:        "obj1",
		Description: "New description",
		Parameters:  params,
	})
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

	pApi.Create(&api.ProjectInfo{Code: ProjectCode, Status: domain.ProjectStatusActive})
	envApi.Create(&api.EnvironmentInfo{Code: envCode})

	objApi.Create(&api.ObjectInfo{Code: "obj1"})

	_, err = objApi.Create(&api.ObjectInfo{
		Code: "obj2",
		Inherits: &domain.ObjectInheritance{
			ProjectCode: ProjectCode,
			EnvCode:     envCode,
			ObjectCode:  "obj1",
		},
	})
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

	pApi.Create(&api.ProjectInfo{Code: ProjectCode, Status: domain.ProjectStatusActive})
	envApi.Create(&api.EnvironmentInfo{Code: envCode})

	objApi.Create(&api.ObjectInfo{
		Code: "obj1",
		Parameters: []*domain.Parameter{
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
		},
	})

	objApi.Create(&api.ObjectInfo{
		Code: "obj2",
		Inherits: &domain.ObjectInheritance{
			ProjectCode: ProjectCode,
			EnvCode:     envCode,
			ObjectCode:  "obj1",
		},
		Parameters: []*domain.Parameter{
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
		},
	})

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

	_, err = objApi.Create(&api.ObjectInfo{
		Code:     "obj3",
		Inherits: &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: envCode, ObjectCode: "obj1"},
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:        "param2",
				Description: "Param 2",
				Type:        domain.ParameterString,
				Value:       "wrong type",
			},
		},
	})
	assert.Equal(api.ErrObjectInheritorTypeMismatch, err)

	obj, err = objApi.Create(&api.ObjectInfo{
		Code:     "obj3",
		Inherits: &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: envCode, ObjectCode: "obj2"},
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:        "param3",
				Description: "Param 3 will be overridden",
				Type:        domain.ParameterString,
				Value:       "value 2",
			},
		},
	})
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

	ls, err := objApi.InheritorsFlatList("obj1")
	assert.Nil(err)
	assert.Len(ls, 2)

	for _, l := range ls {
		switch l.Code {
		case "obj2":
			assert.Equal(domain.ObjectCode("obj1"), l.Inherits.ObjectCode)
			assert.Equal(envCode, l.Inherits.EnvCode)
			assert.Equal(ProjectCode, l.Inherits.ProjectCode)
		case "obj3":
			assert.Equal(domain.ObjectCode("obj2"), l.Inherits.ObjectCode)
			assert.Equal(envCode, l.Inherits.EnvCode)
			assert.Equal(ProjectCode, l.Inherits.ProjectCode)
		}
	}

	AfterTest()

}

func TestObjectsParameterInheritance(t *testing.T) {
	assert := asserts.New(t)
	BeforeTest()

	pApi := GetApi()
	envApi := pApi.For(ProjectCode).Environments()
	objApi := envApi.For(envCode).Objects()

	pApi.Create(&api.ProjectInfo{Code: ProjectCode, Status: domain.ProjectStatusActive})
	envApi.Create(&api.EnvironmentInfo{Code: envCode})

	_, err := objApi.Create(&api.ObjectInfo{
		Code: "obj1",
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:  "param1",
				Type:  domain.ParameterBool,
				Value: true,
			},
		},
	})
	assert.Nil(err)

	_, err = objApi.Create(&api.ObjectInfo{
		Code:     "obj2",
		Inherits: &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: envCode, ObjectCode: "obj1"},
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:  "param1",
				Type:  domain.ParameterBool,
				Value: false,
			},
		},
	})
	assert.Nil(err)

	_, err = objApi.Create(&api.ObjectInfo{
		Code:     "obj3",
		Inherits: &domain.ObjectInheritance{ProjectCode: ProjectCode, EnvCode: envCode, ObjectCode: "obj2"},
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:  "param2",
				Type:  domain.ParameterString,
				Value: "value",
			},
		},
	})
	assert.Nil(err)

	_, err = objApi.Update(&api.ObjectInfo{
		Code: "obj1",
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:  "param1",
				Type:  domain.ParameterString,
				Value: "value",
			},
		},
	})
	assert.IsType(&api.ErrObjectParameter{}, err)
	assert.Equal("Object parameter `param1` error: Object parameter type changing restricted", err.Error())

	_, err = objApi.Update(&api.ObjectInfo{
		Code:       "obj1",
		Parameters: []*domain.Parameter{},
	})
	assert.Nil(err)

	_, err = objApi.Update(&api.ObjectInfo{
		Code: "obj1",
		Parameters: []*domain.Parameter{&domain.Parameter{
			Code:  "param1",
			Type:  domain.ParameterString,
			Value: "value",
		}},
	})
	assert.IsType(&api.ErrObjectParameter{}, err)
	assert.Equal("Object parameter `param1` error: Object parameter exists in inheritor: p1:env_1:obj2", err.Error())

	_, err = objApi.Update(&api.ObjectInfo{
		Code: "obj1",
		Parameters: []*domain.Parameter{
			&domain.Parameter{
				Code:  "param2",
				Type:  domain.ParameterBool,
				Value: false,
			},
		},
	})
	assert.IsType(&api.ErrObjectParameter{}, err)
	assert.Equal("Object parameter `param2` error: Object parameter exists in inheritor: p1:env_1:obj3", err.Error())

	AfterTest()
}
