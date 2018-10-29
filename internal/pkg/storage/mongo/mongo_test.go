package mongo_test

import (
	"testing"

	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/Toggly/core/internal/pkg/storage/mongo"
	"github.com/stretchr/testify/assert"
)

func TestMongoStorage(t *testing.T) {
	t.Skip("not now")
	assert := assert.New(t)
	var err error
	s, err := mongo.NewMongoStorage("mongodb://127.0.0.1:27017/toggly")
	assert.Nil(err)
	assert.NotNil(s)

	ow := s.ForOwner("ow1")
	assert.NotNil(ow)
	projStorage := ow.Projects()

	a, err := projStorage.List()
	assert.Nil(err)
	assert.Empty(a)

	o, err := projStorage.Get("proj1")
	assert.IsType(err, storage.ErrNotFound)
	assert.Nil(o)

	// proj := &domain.Project{
	// 	Code:        "proj1",
	// 	Description: "Project 1",
	// 	RegDate:     time.Now(),
	// 	Status:      domain.ProjectStatusActive,
	// }

	// err = projStorage.Save(proj)

	// assert.Nil(err)
}
