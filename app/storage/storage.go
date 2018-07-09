package storage

import "github.com/Toggly/backend/app/data"

// DataStorage defines storage interface for dictionary
type DataStorage interface {
	ListProjects() ([]*data.Project, error)
	GetProject(code data.CodeType) (*data.Project, error)
	SaveProject(project data.Project) error

	ListEnvironments(project data.CodeType) ([]*data.Environment, error)
	GetEnvironment(project data.CodeType, code data.CodeType) (*data.Environment, error)
	SaveEnvironment(project data.CodeType, env data.Environment) error

	ListObjects(project data.CodeType) ([]*data.Object, error)
	GetObject(project data.CodeType, code data.CodeType) (*data.Object, error)
	SaveObject(project data.CodeType, object data.Object) error
	GetObjectValue(project data.CodeType, code data.CodeType, env data.CodeType, id data.ObjectID) (*data.Object, error)
	SaveObjectValue(project data.CodeType, object data.Object, env data.CodeType, id data.ObjectID) error

	ListParameters(project data.CodeType, object data.CodeType) ([]*data.Parameter, error)
	GetParameter(project data.CodeType, object data.CodeType, code data.CodeType) (*data.Parameter, error)
	SaveParameter(project data.CodeType, object data.CodeType, parameter data.Parameter) error
}
