package api

import (
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
)

// ObjectAPI servers object api namespace
type ObjectAPI struct {
	Storage storage.DataStorage
}

//List returns list of objects
func (o *ObjectAPI) List(owner string, project domain.ProjectCode, env domain.EnvironmentCode) (objects []*domain.Object, err error) {
	objects, err = o.Storage.Projects(owner).For(project).Environments().For("").Objects().List()
	_ = objects
	return objects, err
}
