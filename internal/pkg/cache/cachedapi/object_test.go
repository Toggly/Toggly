package cachedapi_test

import (
	"encoding/json"
	"testing"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/domain"
	asserts "github.com/stretchr/testify/assert"
)

func TestObjectCaching(t *testing.T) {
	assert := asserts.New(t)

	BeforeTest()

	engine, cache := getEngineAndCache()
	engine.ForOwner("ow1").Projects().Create(&api.ProjectInfo{Code: "project1"})
	engine.ForOwner("ow1").Projects().For("project1").Environments().Create(&api.EnvironmentInfo{Code: "env1"})
	eng := engine.ForOwner("ow1").Projects().For("project1").Environments().For("env1").Objects()

	t.Run("empty cache", func(t *testing.T) {
		b, err := cache.Get("/own/ow1/project/project1/env/env1/object")
		assert.Nil(err)
		assert.Nil(b)
	})

	t.Run("create object empty cache", func(t *testing.T) {
		_, err := eng.Create(&api.ObjectInfo{Code: "obj1", Description: "Description 1"})
		assert.Nil(err)
		b, err := cache.Get("/own/ow1/project/project1/env/env1/object")
		assert.Nil(err)
		assert.Nil(b)
	})

	t.Run("list in cache", func(t *testing.T) {
		eng.List()
		var list []*domain.Object
		b, err := cache.Get("/own/ow1/project/project1/env/env1/object")
		assert.Nil(err)
		assert.NotNil(b)
		json.Unmarshal(b, &list)
		assert.Len(list, 1)
	})

	t.Run("object not in cache", func(t *testing.T) {
		b, err := cache.Get("/own/ow1/project/project1/env/env1/object/obj1")
		assert.Nil(err)
		assert.Nil(b)
	})

	_, err := eng.List()
	assert.Nil(err)

	t.Run("object from cache", func(t *testing.T) {
		obj := &domain.Object{}
		eng.Get("obj1")
		b, err := cache.Get("/own/ow1/project/project1/env/env1/object/obj1")
		assert.Nil(err)
		assert.NotNil(b)
		json.Unmarshal(b, obj)
		assert.Equal(domain.ObjectCode("obj1"), obj.Code)
		assert.Equal("Description 1", obj.Description)
		eng.Delete("obj1")
	})

	t.Run("invalidate on update", func(t *testing.T) {
		// Create objects
		_, err := eng.Create(&api.ObjectInfo{Code: "obj1"})
		assert.Nil(err)
		_, err = eng.Create(&api.ObjectInfo{Code: "obj2", Inherits: &domain.ObjectInheritance{ProjectCode: "project1", EnvCode: "env1", ObjectCode: "obj1"}})
		assert.Nil(err)
		_, err = eng.Create(&api.ObjectInfo{Code: "obj3", Inherits: &domain.ObjectInheritance{ProjectCode: "project1", EnvCode: "env1", ObjectCode: "obj2"}})
		assert.Nil(err)
		// Warmup cache
		eng.List()
		eng.Get("obj1")
		eng.Get("obj2")
		eng.Get("obj3")
		b, err := cache.Get("/own/ow1/project/project1/env/env1/object")
		assert.Nil(err)
		assert.NotNil(b)
		b, err = cache.Get("/own/ow1/project/project1/env/env1/object/obj1")
		assert.Nil(err)
		assert.NotNil(b)
		b, err = cache.Get("/own/ow1/project/project1/env/env1/object/obj2")
		assert.Nil(err)
		assert.NotNil(b)
		b, err = cache.Get("/own/ow1/project/project1/env/env1/object/obj3")
		assert.Nil(err)
		assert.NotNil(b)
		// Update
		eng.Update(&api.ObjectInfo{Code: "obj1", Description: "Description 2"})
		b, err = cache.Get("/own/ow1/project/project1/env/object")
		assert.Nil(err)
		assert.Nil(b)
		b, err = cache.Get("/own/ow1/project/project1/env/env1/object/obj1")
		assert.Nil(err)
		assert.Nil(b)
		b, err = cache.Get("/own/ow1/project/project1/env/env1/object/obj2")
		assert.Nil(err)
		assert.Nil(b)
		b, err = cache.Get("/own/ow1/project/project1/env/env1/object/obj3")
		assert.Nil(err)
		assert.Nil(b)
		eng.Delete("obj3")
		eng.Delete("obj2")
		eng.Delete("obj1")
	})

	t.Run("invalidate on delete", func(t *testing.T) {
		t.Skip()
		_, err := eng.Create(&api.ObjectInfo{Code: "obj1", Description: "Description 1"})
		assert.Nil(err)
		eng.List()
		eng.Get("obj1")
		b, err := cache.Get("/own/ow1/project/project1/env/env1/object")
		assert.Nil(err)
		assert.NotNil(b)
		b, err = cache.Get("/own/ow1/project/project1/env/env1/object/obj1")
		assert.Nil(err)
		assert.NotNil(b)
		eng.Delete("obj1")
		b, err = cache.Get("/own/ow1/project/project1/env/env1/object")
		assert.Nil(err)
		assert.Nil(b)
		b, err = cache.Get("/own/ow1/project/project1/env/env1/object/obj1")
		assert.Nil(err)
		assert.Nil(b)
	})

	AfterTest()
}
