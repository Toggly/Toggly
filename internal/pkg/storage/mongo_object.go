package storage

import (
	"errors"

	"github.com/Toggly/core/internal/domain"
)

type mgForEnvironment struct{}

func (s *mgForEnvironment) Objects() ObjectStorage {
	return &mgObjectStorage{}
}

type mgObjectStorage struct{}

func (s *mgObjectStorage) List() ([]*domain.Object, error) {
	return nil, nil
}
func (s *mgObjectStorage) Get(code domain.ObjectCode) (*domain.Object, error) {
	return nil, nil
}
func (s *mgObjectStorage) Delete(code domain.ObjectCode) error {
	return errors.New("Method not implemented")
}
func (s *mgObjectStorage) Save(object domain.Object) error {
	return nil
}
