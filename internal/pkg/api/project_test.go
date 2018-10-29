package api_test

import (
	"testing"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/pkg/api"

	"github.com/Toggly/core/internal/pkg/storage/mongo"
	asserts "github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	assert := asserts.New(t)

	dataStorage, err := mongo.NewMongoStorage("mongodb://localhost:27017/toggly_test")
	assert.Nil(err)

	engine := &api.Engine{Storage: &dataStorage}

	const ow = "test_owner"

	// proj1 := &domain.Project{
	// 	Code:        "project1",
	// 	Description: "Project 1 Description",
	// 	Status:      domain.ProjectStatusActive,
	// 	RegDate:     ,
	// }

	pApi := engine.ForOwner(ow).Projects()

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

	assert.Nil(pApi.Delete("p1"))
	assert.Equal(pApi.Delete("p1"), api.ErrProjectNotFound)
	assert.Equal(pApi.Delete("not_existing_code"), api.ErrProjectNotFound)
}
