package api

import (
	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/pkg/storage"
)

// ProjectAPI servers project api namespace
type ProjectAPI struct {
	Storage storage.DataStorage
}

//List returns list of projects
func (p *ProjectAPI) List(owner string) ([]*domain.Project, error) {
	return p.Storage.Projects(owner).List()
}

// Get Project By ID
func (p *ProjectAPI) Get(owner string, id string) (*domain.Project, error) {
	return p.Storage.Projects(owner).Get(domain.ProjectCode(id))
}

// Save Project
func (p *ProjectAPI) Save(owner string, project *domain.Project) error {
	return p.Storage.Projects(owner).Save(*project)
}
