package engine

import (
	"time"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/globalsign/mgo/bson"
)

// ProjectAPI servers project api namespace
type ProjectAPI struct {
	OwnerAPI
}

func (p *ProjectAPI) storage() storage.ProjectStorage {
	return (*p.Storage).ForOwner(p.Owner).Projects()
}

//List returns list of projects
func (p *ProjectAPI) List() ([]*domain.Project, error) {
	return p.storage().List()
}

// Get Project By code
func (p *ProjectAPI) Get(code domain.ProjectCode) (*domain.Project, error) {
	project, err := p.storage().Get(code)
	if err == storage.ErrNotFound {
		return nil, api.ErrProjectNotFound
	}
	return project, err
}

func checkProjectParams(code domain.ProjectCode, description string, status domain.ProjectStatus) error {
	if code == "" {
		return api.NewBadRequestError("Project code not specified")
	}
	return nil
}

// Create Project
func (p *ProjectAPI) Create(info *api.ProjectInfo) (*domain.Project, error) {
	code := info.Code
	description := info.Description
	status := info.Status
	if err := checkProjectParams(code, description, status); err != nil {
		return nil, err
	}
	newProj := &domain.Project{
		OwnerID:     p.Owner,
		Code:        code,
		Description: description,
		RegDate:     bson.Now().In(time.UTC),
		Status:      status,
	}
	if err := p.storage().Save(newProj); err != nil {
		return nil, err
	}
	return newProj, nil
}

// Update Project
func (p *ProjectAPI) Update(info *api.ProjectInfo) (*domain.Project, error) {
	code := info.Code
	description := info.Description
	status := info.Status
	if err := checkProjectParams(code, description, status); err != nil {
		return nil, err
	}
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
	if err := p.storage().Update(newProj); err != nil {
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
		return api.ErrProjectNotEmpty
	}
	err = p.storage().Delete(code)
	if err == storage.ErrNotFound {
		return api.ErrProjectNotFound
	}
	return err
}

// For returns environment api for specified project
func (p *ProjectAPI) For(code domain.ProjectCode) api.ForProjectAPI {
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
func (fp *ForProjectAPI) Environments() api.EnvironmentAPI {
	return &EnvironmentAPI{
		Owner:       fp.Owner,
		ProjectCode: fp.ProjectCode,
		Storage:     fp.Storage,
		ProjectAPI:  fp.ProjectAPI,
	}
}
