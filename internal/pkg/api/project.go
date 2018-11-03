package api

import (
	"time"

	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/globalsign/mgo/bson"
)

// ProjectAPI servers project api namespace
type ProjectAPI struct {
	OwnerAPI
}

//List returns list of projects
func (p *ProjectAPI) List() ([]*domain.Project, error) {
	return (*p.Storage).ForOwner(p.Owner).Projects().List()
}

// Get Project By code
func (p *ProjectAPI) Get(code domain.ProjectCode) (*domain.Project, error) {
	project, err := (*p.Storage).ForOwner(p.Owner).Projects().Get(code)
	if err == storage.ErrNotFound {
		return nil, ErrProjectNotFound
	}
	return project, err
}

// Create Project
func (p *ProjectAPI) Create(code domain.ProjectCode, description string, status domain.ProjectStatus) (*domain.Project, error) {
	newProj := &domain.Project{
		OwnerID:     p.Owner,
		Code:        code,
		Description: description,
		RegDate:     bson.Now().In(time.UTC),
		Status:      status,
	}
	if err := (*p.Storage).ForOwner(p.Owner).Projects().Save(newProj); err != nil {
		return nil, err
	}
	return newProj, nil
}

// Update Project
func (p *ProjectAPI) Update(code domain.ProjectCode, description string, status domain.ProjectStatus) (*domain.Project, error) {
	pr, err := p.Get(code)
	if err != nil {
		return nil, err
	}
	newProj := &domain.Project{
		OwnerID:     p.Owner,
		Code:        code,
		Description: description,
		RegDate:     pr.RegDate,
		Status:      status,
	}
	if err := (*p.Storage).ForOwner(p.Owner).Projects().Update(newProj); err != nil {
		return nil, err
	}
	return newProj, nil
}

// Delete Project
func (p *ProjectAPI) Delete(code domain.ProjectCode) error {
	envList, err := p.For(code).Environments().List()
	if err != nil {
		return err
	}
	if len(envList) > 0 {
		return ErrProjectNotEmpty
	}
	err = (*p.Storage).ForOwner(p.Owner).Projects().Delete(code)
	if err == storage.ErrNotFound {
		return ErrProjectNotFound
	}
	return err
}

// For returns environment api for specified project
func (p *ProjectAPI) For(code domain.ProjectCode) *ForProjectAPI {
	return &ForProjectAPI{
		Owner:       p.Owner,
		ProjectCode: code,
		Storage:     p.Storage,
		ProjectAPI:  p,
	}
}

// ForProjectAPI type
type ForProjectAPI struct {
	Owner       string
	ProjectCode domain.ProjectCode
	Storage     *storage.DataStorage
	ProjectAPI  *ProjectAPI
}

// Environments returns Environments API
func (fp *ForProjectAPI) Environments() *EnvironmentAPI {
	return &EnvironmentAPI{
		Owner:       fp.Owner,
		ProjectCode: fp.ProjectCode,
		Storage:     fp.Storage,
		ProjectAPI:  fp.ProjectAPI,
	}
}
