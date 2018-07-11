package storage

import "github.com/Toggly/backend/app/data"

// NewMongo returns MongoDB storage implementation
func NewMongo() (DataStorage, error) {
	return &mongoStorage{}, nil
}

type mongoStorage struct{}

func (s *mongoStorage) ListProjects() ([]*data.Project, error) {
	return nil, nil
}

func (s *mongoStorage) GetProject(code data.CodeType) (*data.Project, error) {
	return &data.Project{}, nil
}

func (s *mongoStorage) SaveProject(project data.Project) error {
	return nil
}

func (s *mongoStorage) ListEnvironments(project data.CodeType) ([]*data.Environment, error) {
	environments := make([]*data.Environment, 0)
	environments = append(environments, &data.Environment{})
	return environments, nil
}

func (s *mongoStorage) GetEnvironment(project data.CodeType, code data.CodeType) (*data.Environment, error) {
	return &data.Environment{}, nil
}
func (s *mongoStorage) SaveEnvironment(project data.CodeType, env data.Environment) error {
	return nil
}

func (s *mongoStorage) ListObjects(project data.CodeType) ([]*data.Object, error) {
	objects := make([]*data.Object, 0)
	objects = append(objects, &data.Object{})
	return objects, nil
}

func (s *mongoStorage) GetObject(project data.CodeType, code data.CodeType) (*data.Object, error) {
	return &data.Object{}, nil
}
func (s *mongoStorage) SaveObject(project data.CodeType, object data.Object) error {
	return nil
}

func (s *mongoStorage) GetObjectValue(project data.CodeType, code data.CodeType, env data.CodeType, id data.ObjectID) (*data.Object, error) {
	return &data.Object{}, nil
}

func (s *mongoStorage) SaveObjectValue(project data.CodeType, object data.Object, env data.CodeType, id data.ObjectID) error {
	return nil
}

func (s *mongoStorage) ListParameters(project data.CodeType, object data.CodeType) ([]*data.Parameter, error) {
	parameters := make([]*data.Parameter, 0)
	parameters = append(parameters, &data.Parameter{})
	return parameters, nil
}

func (s *mongoStorage) GetParameter(project data.CodeType, object data.CodeType, code data.CodeType) (*data.Parameter, error) {
	return &data.Parameter{}, nil
}

func (s *mongoStorage) SaveParameter(project data.CodeType, object data.CodeType, parameter data.Parameter) error {
	return nil
}
