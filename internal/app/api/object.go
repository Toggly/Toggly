package api

import (
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
)

// ObjectAPI servers object api namespace
type ObjectAPI struct {
	Owner   string
	Project domain.ProjectCode
	Env     domain.EnvironmentCode
	Storage *storage.DataStorage
}

//List returns list of objects
func (o *ObjectAPI) List() (objects []*domain.Object, err error) {
	objects, err = (*o.Storage).ForOwner(o.Owner).Projects().For(o.Project).Environments().For("").Objects().List()
	_ = objects
	return objects, err
}
