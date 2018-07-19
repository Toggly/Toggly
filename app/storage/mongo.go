package storage

import "github.com/Toggly/core/app/data"

// NewMongoStorage implements DataStorage interface for MongoDB
func NewMongoStorage() (DataStorage, error) {
	return &mgStorage{}, nil
}

type mgStorage struct{}

func (s *mgStorage) Projects() ProjectStorage {
	return &mgProjectStorage{}
}

type mgProjectStorage struct{}

func (s *mgProjectStorage) List() ([]*data.Project, error) {
	return nil, nil
}

func (s *mgProjectStorage) Get(code data.ProjectCode) (*data.Project, error) {
	return nil, nil
}

func (s *mgProjectStorage) Save(project data.Project) error {
	return nil
}

func (s *mgProjectStorage) For(project data.ProjectCode) ForProject {
	return &mgForProject{
		project: project,
	}
}

type mgForProject struct {
	project data.ProjectCode
}

func (s *mgForProject) Environments() EnvironmentStorage {
	return &mgEnvironmentStorage{
		project: s.project,
	}
}

type mgEnvironmentStorage struct {
	project data.ProjectCode
}

func (s *mgEnvironmentStorage) List() ([]*data.Environment, error) {
	return nil, nil
}

func (s *mgEnvironmentStorage) Get(code data.EnvironmentCode) (*data.Environment, error) {
	return nil, nil
}

func (s *mgEnvironmentStorage) Save(env data.Environment) error {
	return nil
}

func (s *mgEnvironmentStorage) For(code data.EnvironmentCode) ForEnvironment {
	return &mgForEnvironment{}
}

type mgForEnvironment struct{}

func (s *mgForEnvironment) Objects() ObjectStorage {
	return &mgObjectStorage{}
}

type mgObjectStorage struct{}

func (s *mgObjectStorage) List() ([]*data.Object, error) {
	return nil, nil
}
func (s *mgObjectStorage) Get(code data.ObjectCode) (*data.Object, error) {
	return nil, nil
}
func (s *mgObjectStorage) Save(object data.Object) error {
	return nil
}
