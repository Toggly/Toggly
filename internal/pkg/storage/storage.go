package storage

import (
	"fmt"

	"github.com/Toggly/core/internal/domain"
)

// UniqueIndexError type
type UniqueIndexError struct {
	Type string
	Key  string
}

func (e *UniqueIndexError) Error() string {
	return fmt.Sprintf("Unique index error: %s [%s]", e.Type, e.Key)
}

// DataStorage defines storage interface
type DataStorage interface {
	ForOwner(ownerID string) OwnerStorage
}

// OwnerStorage defines owner storage interface
type OwnerStorage interface {
	Projects() ProjectStorage
}

// ProjectStorage defines projects storage interface
type ProjectStorage interface {
	List() ([]*domain.Project, error)
	Get(code domain.ProjectCode) (*domain.Project, error)
	Save(project domain.Project) error
	For(project domain.ProjectCode) ForProject
}

// ForProject defines project dependencies interface
type ForProject interface {
	Environments() EnvironmentStorage
}

// EnvironmentStorage defines environment storage interface
type EnvironmentStorage interface {
	List() ([]*domain.Environment, error)
	Get(code domain.EnvironmentCode) (*domain.Environment, error)
	Save(env domain.Environment) error
	For(domain.EnvironmentCode) ForEnvironment
}

// ForEnvironment defines environment dependencies interface
type ForEnvironment interface {
	Objects() ObjectStorage
}

// ObjectStorage defines object structure storage interface
type ObjectStorage interface {
	List() ([]*domain.Object, error)
	Get(code domain.ObjectCode) (*domain.Object, error)
	Save(object domain.Object) error
}
