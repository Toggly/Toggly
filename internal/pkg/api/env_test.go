package api_test

import (
	"testing"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/pkg/api"

	"github.com/Toggly/core/internal/pkg/storage/mongo"
	asserts "github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	assert := asserts.New(t)

	dataStorage, err := mongo.NewMongoStorage("mongodb://localhost:27017/toggly_test")
	assert.Nil(err)

	engine := &api.Engine{Storage: &dataStorage}

	const ow = "test_owner"

	pApi := engine.ForOwner(ow).Projects()
	pApi.Delete("p10")

	pl, err := pApi.List()
	assert.Nil(err)
	assert.Len(pl, 0)

	pr, err := pApi.Create("p10", "Description 10", domain.ProjectStatusActive)
	assert.Nil(err)
	assert.NotNil(pr)
	assert.Equal(domain.ProjectCode("p10"), pr.Code)
	assert.Equal("Description 1", pr.Description)
	assert.Equal(ow, pr.OwnerID)
	assert.NotNil(pr.RegDate)
	assert.Equal(domain.ProjectStatusActive, pr.Status)

	assert.Nil(pApi.Delete("p10"))
}
