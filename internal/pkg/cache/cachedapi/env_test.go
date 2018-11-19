package cachedapi_test

import (
	"encoding/json"
	"testing"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/domain"
	asserts "github.com/stretchr/testify/assert"
)

func TestEnvCaching(t *testing.T) {
	assert := asserts.New(t)

	BeforeTest()

	engine, cache := getEngineAndCache()
	engine.ForOwner("ow1").Projects().Create(&api.ProjectInfo{Code: "project1"})
	eng := engine.ForOwner("ow1").Projects().For("project1").Environments()

	b, err := cache.Get("/own/ow1/project/project1/env")
	assert.Nil(err)
	assert.Nil(b)

	eng.Create(&api.EnvironmentInfo{Code: "env1", Description: "Description 1"})

	b, err = cache.Get("/own/ow1/project/project1/env")
	assert.Nil(err)
	assert.Nil(b)

	eng.List()

	var list []*domain.Environment

	b, err = cache.Get("/own/ow1/project/project1/env")
	assert.Nil(err)
	assert.NotNil(b)
	json.Unmarshal(b, &list)
	assert.Len(list, 1)

	b, err = cache.Get("/own/ow1/project/project1/env/env1")
	assert.Nil(err)
	assert.Nil(b)

	_, err = eng.List()
	assert.Nil(err)

	env := &domain.Environment{}

	eng.Get("env1")
	b, err = cache.Get("/own/ow1/project/project1/env/env1")
	assert.Nil(err)
	assert.NotNil(b)
	json.Unmarshal(b, env)
	assert.Equal(domain.EnvironmentCode("env1"), env.Code)
	assert.Equal("Description 1", env.Description)

	eng.Update(&api.EnvironmentInfo{Code: "env1", Description: "Description 2"})
	b, err = cache.Get("/own/ow1/project/project1/env")
	assert.Nil(err)
	assert.Nil(b)
	b, err = cache.Get("/own/ow1/project/project1/env/env1")
	assert.Nil(err)
	assert.Nil(b)

	eng.Create(&api.EnvironmentInfo{Code: "env1", Description: "Description 1"})
	eng.List()
	eng.Get("env1")
	b, err = cache.Get("/own/ow1/project/project1/env")
	assert.Nil(err)
	assert.NotNil(b)
	b, err = cache.Get("/own/ow1/project/project1/env/env1")
	assert.Nil(err)
	assert.NotNil(b)

	eng.Delete("env1")
	b, err = cache.Get("/own/ow1/project/project1/env")
	assert.Nil(err)
	assert.Nil(b)
	b, err = cache.Get("/own/ow1/project/project1/env/env1")
	assert.Nil(err)
	assert.Nil(b)

	AfterTest()
}
