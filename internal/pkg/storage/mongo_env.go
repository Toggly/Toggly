package storage

import (
	"errors"

	"github.com/Toggly/core/internal/domain"
)

type mgForProject struct {
	project domain.ProjectCode
}

func (s *mgForProject) Environments() EnvironmentStorage {
	return &mgEnvironmentStorage{
		project: s.project,
	}
}

type mgEnvironmentStorage struct {
	project domain.ProjectCode
}

func (s *mgEnvironmentStorage) List() ([]*domain.Environment, error) {
	return nil, nil
}

func (s *mgEnvironmentStorage) Get(code domain.EnvironmentCode) (*domain.Environment, error) {
	return nil, nil
}

func (s *mgEnvironmentStorage) Delete(code domain.EnvironmentCode) error {
	return errors.New("Method not implemented")
}

func (s *mgEnvironmentStorage) Save(env domain.Environment) error {
	return nil
}

func (s *mgEnvironmentStorage) For(code domain.EnvironmentCode) ForEnvironment {
	return &mgForEnvironment{}
}
