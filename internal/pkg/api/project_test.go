package api_test

import (
	"testing"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/pkg/api"
	"github.com/Toggly/core/internal/pkg/storage"

	"github.com/Toggly/core/internal/pkg/storage/mongo"
	asserts "github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	assert := asserts.New(t)

	dataStorage, err := mongo.NewMongoStorage("mongodb://localhost:27017/toggly_test")
	assert.Nil(err)

	engine := &api.Engine{Storage: &dataStorage}

	const ow = "test_owner"

	pApi := engine.ForOwner(ow).Projects()

	pApi.Delete("p1")

	pl, err := pApi.List()
	assert.Nil(err)
	assert.Len(pl, 0)

	pr, err := pApi.Get("p1")
	assert.Equal(err, api.ErrProjectNotFound)
	assert.Nil(pr)

	pr, err = pApi.Create("p1", "Description 1", domain.ProjectStatusActive)
	assert.Nil(err)
	assert.NotNil(pr)
	assert.Equal(domain.ProjectCode("p1"), pr.Code)
	assert.Equal("Description 1", pr.Description)
	assert.Equal(ow, pr.OwnerID)
	assert.NotNil(pr.RegDate)
	assert.Equal(domain.ProjectStatusActive, pr.Status)

	pl, err = pApi.List()
	assert.Len(pl, 1)

	_, err = pApi.Create("p1", "Description 1", domain.ProjectStatusActive)
	assert.NotNil(err)
	assert.IsType(&storage.UniqueIndexError{}, err)

	pr1, err := pApi.Get("p1")
	assert.Nil(err)
	assert.NotNil(pr1)
	assert.Equal(domain.ProjectCode("p1"), pr1.Code)
	assert.Equal("Description 1", pr1.Description)
	assert.Equal(ow, pr1.OwnerID)
	assert.Equal(pr.RegDate, pr1.RegDate)
	assert.Equal(pr.Status, pr1.Status)

	pr1u, err := pApi.Update("p1", "Description 2", domain.ProjectStatusDisabled)
	assert.Nil(err)
	assert.NotNil(pr1u)
	assert.Equal(domain.ProjectCode("p1"), pr1u.Code)
	assert.Equal("Description 2", pr1u.Description)
	assert.Equal(ow, pr1u.OwnerID)
	assert.Equal(pr.RegDate, pr1u.RegDate)
	assert.Equal(domain.ProjectStatusDisabled, pr1u.Status)

	pr2u, err := pApi.Update("p2", "Description 2", domain.ProjectStatusDisabled)
	assert.Nil(pr2u)
	assert.Equal(api.ErrProjectNotFound, err)

	pApi.For("p1").Environments().Create("env_code", "", false)

	assert.Equal(api.ErrProjectNotEmpty, pApi.Delete("p1"))

	assert.Nil(pApi.For("p1").Environments().Delete("env_code"))
	assert.Nil(pApi.Delete("p1"))

	assert.Equal(api.ErrProjectNotFound, pApi.Delete("p1"))
}
