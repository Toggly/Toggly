package storage

import (
	"time"

	"github.com/Toggly/backend/app/data"
)

// NewFake returns fake storage implementation
func NewFake() (DataStorage, error) {
	return &fakeStorage{}, nil
}

type fakeStorage struct{}

func (s *fakeStorage) ListProjects() ([]*data.Project, error) {
	return projects(), nil
}

func (s *fakeStorage) GetProject(code data.CodeType) (*data.Project, error) {
	return findProj(string(code)), nil
}

func (s *fakeStorage) SaveProject(project data.Project) error {
	return nil
}

func (s *fakeStorage) ListEnvironments(project data.CodeType) ([]*data.Environment, error) {
	return envs(), nil
}

func (s *fakeStorage) GetEnvironment(project data.CodeType, code data.CodeType) (*data.Environment, error) {
	return findEnv(string(code)), nil
}
func (s *fakeStorage) SaveEnvironment(project data.CodeType, env data.Environment) error {
	return nil
}

func (s *fakeStorage) ListObjects(project data.CodeType) ([]*data.Object, error) {
	return objects(), nil
}

func (s *fakeStorage) GetObject(project data.CodeType, code data.CodeType) (*data.Object, error) {
	return findObj(string(code)), nil
}
func (s *fakeStorage) SaveObject(project data.CodeType, object data.Object) error {
	return nil
}

func (s *fakeStorage) GetObjectValue(project data.CodeType, code data.CodeType, env data.CodeType, id data.ObjectID) (*data.Object, error) {
	return &data.Object{}, nil
}

func (s *fakeStorage) SaveObjectValue(project data.CodeType, object data.Object, env data.CodeType, id data.ObjectID) error {
	return nil
}

func (s *fakeStorage) ListParameters(project data.CodeType, object data.CodeType) ([]*data.Parameter, error) {
	parameters := make([]*data.Parameter, 0)
	parameters = append(parameters, &data.Parameter{})
	return parameters, nil
}

func (s *fakeStorage) GetParameter(project data.CodeType, object data.CodeType, code data.CodeType) (*data.Parameter, error) {
	return &data.Parameter{}, nil
}

func (s *fakeStorage) SaveParameter(project data.CodeType, object data.CodeType, parameter data.Parameter) error {
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

func findObj(code string) *data.Object {
	for _, v := range objects() {
		if string(v.Code) == code {
			return v
		}
	}
	return nil
}

func projects() []*data.Project {
	p := make([]*data.Project, 3)

	p[0] = &data.Project{
		Code:        "project1",
		Description: "Simple Project 1",
		RegDate:     time.Now(),
		Status:      0,
	}
	p[1] = &data.Project{
		Code:        "project2",
		Description: "Simple Project 2",
		RegDate:     time.Now(),
		Status:      0,
	}
	p[2] = &data.Project{
		Code:        "project3",
		Description: "Simple Project 3",
		RegDate:     time.Now(),
		Status:      1,
	}
	return p
}

func findProj(code string) *data.Project {
	for _, v := range projects() {
		if string(v.Code) == code {
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

func findEnv(code string) *data.Environment {
	for _, v := range envs() {
		if string(v.Code) == code {
			return v
		}
	}
	return nil
}
