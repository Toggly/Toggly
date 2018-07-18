package storage

import (
	"time"

	"github.com/Toggly/core/app/data"
)

// NewHashMapStorage returns hashmap storage implementation
func NewHashMapStorage() (DataStorage, error) {
	return &hmStorage{}, nil
}

type hmStorage struct{}

func (s *hmStorage) Projects() ProjectStorage {
	return &hmProjectStorage{}
}

type hmProjectStorage struct{}

func (s *hmProjectStorage) List() ([]*data.Project, error) {
	return projects(), nil
}

func (s *hmProjectStorage) Get(code data.ProjectCode) (*data.Project, error) {
	return findProj(code), nil
}

func (s *hmProjectStorage) Save(project data.Project) error {
	return nil
}

func (s *hmProjectStorage) For(project data.ProjectCode) ForProject {
	return &hmForProject{
		project: project,
	}
}

type hmForProject struct {
	project data.ProjectCode
}

func (s *hmForProject) Environments() EnvironmentStorage {
	return &hmEnvironmentStorage{
		project: s.project,
	}
}

type hmEnvironmentStorage struct {
	project data.ProjectCode
}

func (s *hmEnvironmentStorage) List() ([]*data.Environment, error) {
	return envs(), nil
}

func (s *hmEnvironmentStorage) Get(code data.EnvironmentCode) (*data.Environment, error) {
	return findEnv(code), nil
}

func (s *hmEnvironmentStorage) Save(env data.Environment) error {
	return nil
}

func (s *hmEnvironmentStorage) For(code data.EnvironmentCode) ForEnvironment {
	return &hmForEnvironment{}
}

type hmForEnvironment struct{}

func (s *hmForEnvironment) Objects() ObjectStorage {
	return &hmObjectStorage{}
}

type hmObjectStorage struct{}

func (s *hmObjectStorage) List() ([]*data.Object, error) {
	return objects(), nil
}
func (s *hmObjectStorage) Get(code data.ObjectCode) (*data.Object, error) {
	return findObj(code), nil
}
func (s *hmObjectStorage) Save(object data.Object) error {
	return nil
}

// FAKE DATA

func objects() []*data.Object {
	p := make([]*data.Object, 2)
	p[0] = &data.Object{
		Code:        "user",
		Description: "User object",
		Inherits:    "group",
		Parameters:  make([]data.Parameter, 0),
	}
	p[1] = &data.Object{
		Code:        "group",
		Description: "Group object",
		Inherits:    "",
		Parameters:  make([]data.Parameter, 0),
	}
	return p
}

func findObj(code data.ObjectCode) *data.Object {
	for _, v := range objects() {
		if v.Code == code {
			return v
		}
	}
	return nil
}

func projects() []*data.Project {
	p := make([]*data.Project, 3)

	p[0] = &data.Project{
		ID:          "0",
		Code:        "project1",
		Description: "Simple Project 1",
		RegDate:     time.Now(),
		Status:      0,
	}
	p[1] = &data.Project{
		ID:          "1",
		Code:        "project2",
		Description: "Simple Project 2",
		RegDate:     time.Now(),
		Status:      0,
	}
	p[2] = &data.Project{
		ID:          "2",
		Code:        "project3",
		Description: "Simple Project 3",
		RegDate:     time.Now(),
		Status:      1,
	}
	return p
}

func findProj(code data.ProjectCode) *data.Project {
	for _, v := range projects() {
		if v.Code == code {
			return v
		}
	}
	return nil
}

func envs() []*data.Environment {
	p := make([]*data.Environment, 2)
	p[0] = &data.Environment{
		Code:        "dev",
		Description: "Development environment",
		Protected:   false,
	}
	p[1] = &data.Environment{
		Code:        "prod",
		Description: "Production environment",
		Protected:   true,
	}
	return p
}

func findEnv(code data.EnvironmentCode) *data.Environment {
	for _, v := range envs() {
		if v.Code == code {
			return v
		}
	}
	return nil
}
