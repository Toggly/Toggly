package storage

import "github.com/Toggly/core/app/data"

// DataStorage defines storage interface for dictionary
type DataStorage interface {
	Projects() ProjectStorage
}

// ProjectStorage defines projects storage interface
type ProjectStorage interface {
	List() ([]*data.Project, error)
	Get(code data.ProjectCode) (*data.Project, error)
	Save(project data.Project) error
	For(project data.ProjectCode) ForProject
}

// ForProject defines project dependencies interface
type ForProject interface {
	Environments() EnvironmentStorage
}

// EnvironmentStorage defines environment storage interface
type EnvironmentStorage interface {
	List() ([]*data.Environment, error)
	Get(code data.EnvironmentCode) (*data.Environment, error)
	Save(env data.Environment) error
	For(data.EnvironmentCode) ForEnvironment
}

// ForEnvironment defines environment dependencies interface
type ForEnvironment interface {
	Objects() ObjectStorage
}

// ObjectStorage defines object structure storage interface
type ObjectStorage interface {
	List() ([]*data.Object, error)
	Get(code data.ObjectCode) (*data.Object, error)
	Save(object data.Object) error
}
