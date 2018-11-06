package api_test

import (
	"testing"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/pkg/api"
	"github.com/Toggly/core/internal/pkg/storage"

	asserts "github.com/stretchr/testify/assert"
)

const envCode = domain.EnvironmentCode("env_1")

func TestEnvWithNoProject(t *testing.T) {
	assert := asserts.New(t)

	BeforeTest()

	pApi := GetApi()
	envApi := pApi.For(ProjectCode).Environments()

	envList, err := envApi.List()
	assert.Equal(api.ErrProjectNotFound, err)
	assert.Nil(envList)

	env, err := envApi.Get(envCode)
	assert.Equal(api.ErrProjectNotFound, err)
	assert.Nil(env)

	env, err = envApi.Create(envCode, "", false)
	assert.Equal(api.ErrProjectNotFound, err)
	assert.Nil(env)

	env, err = envApi.Update(envCode, "", false)
	assert.Equal(api.ErrProjectNotFound, err)
	assert.Nil(env)

	err = envApi.Delete(envCode)
	assert.Equal(api.ErrProjectNotFound, err)

	AfterTest()
}

func TestEnvWithProject(t *testing.T) {
	assert := asserts.New(t)

	BeforeTest()

	pApi := GetApi()
	envApi := pApi.For(ProjectCode).Environments()

	pApi.Create(ProjectCode, "Description 1", domain.ProjectStatusActive)

	envList, err := envApi.List()
	assert.Nil(err)
	assert.Empty(envList)

	env, err := envApi.Get(envCode)
	assert.Equal(api.ErrEnvironmentNotFound, err)
	assert.Nil(env)

	env, err = envApi.Update(envCode, "", false)
	assert.Equal(api.ErrEnvironmentNotFound, err)
	assert.Nil(env)

	err = envApi.Delete(envCode)
	assert.Equal(api.ErrEnvironmentNotFound, err)

	env, err = envApi.Create(envCode, "Description 1", false)
	assert.Nil(err)
	assert.NotNil(env)
	assert.Equal(ow, env.OwnerID)
	assert.Equal(ProjectCode, env.ProjectCode)
	assert.Equal(envCode, env.Code)
	assert.Equal("Description 1", env.Description)
	assert.False(env.Protected)
	assert.NotNil(env.RegDate)

	_, err = envApi.Create(envCode, "Description 1", false)
	assert.IsType(&storage.UniqueIndexError{}, err)

	envList, err = envApi.List()
	assert.Nil(err)
	assert.Len(envList, 1)

	env, err = envApi.Get(envCode)
	assert.Nil(err)
	assert.NotNil(env)
	assert.Equal(ow, env.OwnerID)
	assert.Equal(ProjectCode, env.ProjectCode)
	assert.Equal(envCode, env.Code)
	assert.Equal("Description 1", env.Description)
	assert.False(env.Protected)
	assert.NotNil(env.RegDate)

	envU, err := envApi.Update(envCode, "Description 2", true)
	assert.Nil(err)
	assert.NotNil(envU)
	assert.Equal(envCode, envU.Code)
	assert.Equal(ow, envU.OwnerID)
	assert.Equal(ProjectCode, envU.ProjectCode)
	assert.Equal("Description 2", envU.Description)
	assert.True(envU.Protected)
	assert.Equal(env.RegDate, envU.RegDate)

	envU2, err := envApi.Update("env_2", "Description 2", true)
	assert.Nil(envU2)
	assert.Equal(api.ErrEnvironmentNotFound, err)

	envApi.For(envCode).Objects().Create(domain.ObjectCode("obj1"), "", nil, nil)

	assert.Equal(api.ErrEnvironmentNotEmpty, envApi.Delete(envCode))

	envApi.For(envCode).Objects().Delete(domain.ObjectCode("obj1"))

	assert.Nil(envApi.Delete(envCode))

	AfterTest()
}
