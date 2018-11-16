package api

import (
	"errors"
	"fmt"

	"github.com/Toggly/core/internal/domain"
)

var (
	// ErrProjectNotFound error
	ErrProjectNotFound = errors.New("Project not found")
	// ErrProjectNotEmpty error
	ErrProjectNotEmpty = errors.New("Project not empty")

	// ErrEnvironmentNotFound error
	ErrEnvironmentNotFound = errors.New("Environment not found")
	// ErrEnvironmentNotEmpty error
	ErrEnvironmentNotEmpty = errors.New("Environment not empty")

	// ErrObjectNotFound error
	ErrObjectNotFound = errors.New("Object not found")
	// ErrObjectHasInheritors error
	ErrObjectHasInheritors = errors.New("Object has inheritors")
	// ErrObjectParentNotExists error
	ErrObjectParentNotExists = errors.New("Object parrent does not exists")
	// ErrObjectInheritorTypeMismatch error
	ErrObjectInheritorTypeMismatch = errors.New("Object inheritor parameter type mismatch")
)

// ErrBadRequest type
type ErrBadRequest struct {
	Description string
}

func (e *ErrBadRequest) Error() string {
	return fmt.Sprintf("Bad request: %s", e.Description)
}

// NewBadRequestError returns ErrBadRequest error
func NewBadRequestError(description string) *ErrBadRequest {
	return &ErrBadRequest{Description: description}
}

// ErrObjectParameter type
type ErrObjectParameter struct {
	Name        string
	Description string
}

func (e *ErrObjectParameter) Error() string {
	return fmt.Sprintf("Object parameter `%s` error: %s", e.Name, e.Description)
}

// NewObjectParameterError returns ErrBadRequest error
func NewObjectParameterError(name, description string) *ErrObjectParameter {
	return &ErrObjectParameter{
		Name:        name,
		Description: description,
	}
}

// TogglyAPI interface
type TogglyAPI interface {
	ForOwner(owner string) OwnerAPI
}

// OwnerAPI interface
type OwnerAPI interface {
	Projects() ProjectAPI
}

// ProjectInfo type
type ProjectInfo struct {
	Code        domain.ProjectCode
	Description string
	Status      domain.ProjectStatus
}

// ProjectAPI interface
type ProjectAPI interface {
	List() ([]*domain.Project, error)
	Get(code domain.ProjectCode) (*domain.Project, error)
	Create(info *ProjectInfo) (*domain.Project, error)
	Update(info *ProjectInfo) (*domain.Project, error)
	Delete(code domain.ProjectCode) error
	For(code domain.ProjectCode) ForProjectAPI
}

// ForProjectAPI interface
type ForProjectAPI interface {
	Environments() EnvironmentAPI
}

// EnvironmentInfo type
type EnvironmentInfo struct {
	Code        domain.EnvironmentCode
	Description string
	Protected   bool
}

// EnvironmentAPI interface
type EnvironmentAPI interface {
	List() ([]*domain.Environment, error)
	Get(code domain.EnvironmentCode) (*domain.Environment, error)
	Create(info *EnvironmentInfo) (*domain.Environment, error)
	Update(info *EnvironmentInfo) (*domain.Environment, error)
	Delete(code domain.EnvironmentCode) error
	For(code domain.EnvironmentCode) ForObjectAPI
}

// ForObjectAPI interface
type ForObjectAPI interface {
	Objects() ObjectAPI
}

// ObjectInfo type
type ObjectInfo struct {
	Code        domain.ObjectCode
	Description string
	Inherits    *domain.ObjectInheritance
	Parameters  []*domain.Parameter
}

// ObjectAPI interface
type ObjectAPI interface {
	List() ([]*domain.Object, error)
	Get(code domain.ObjectCode) (*domain.Object, error)
	Create(info *ObjectInfo) (*domain.Object, error)
	Update(info *ObjectInfo) (*domain.Object, error)
	Delete(code domain.ObjectCode) error
}
