package api

import (
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
)

// ProjectAPI servers project api namespace
type ProjectAPI struct {
	Owner   string
	Storage *storage.DataStorage
}

//List returns list of projects
func (p *ProjectAPI) List() ([]*domain.Project, error) {
	return (*p.Storage).ForOwner(p.Owner).Projects().List()
}

// Get Project By code
func (p *ProjectAPI) Get(code domain.ProjectCode) (*domain.Project, error) {
	return (*p.Storage).ForOwner(p.Owner).Projects().Get(code)
}

// Save Project
func (p *ProjectAPI) Save(project *domain.Project) error {
	return (*p.Storage).ForOwner(p.Owner).Projects().Save(*project)
}

// For returns environment api for specified project
func (p *ProjectAPI) For(code domain.ProjectCode) *EnvironmentAPI {
	return &EnvironmentAPI{
		Owner:       p.Owner,
		ProjectCode: code,
		Storage:     p.Storage,
	}
}
