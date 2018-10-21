package mongo

import (
	"errors"

	"github.com/Toggly/core/internal/domain"
)

type mgoObjectStorage struct{}

func (s *mgoObjectStorage) List() ([]*domain.Object, error) {
	return nil, nil
}
func (s *mgoObjectStorage) Get(code domain.ObjectCode) (*domain.Object, error) {
	return nil, nil
}
func (s *mgoObjectStorage) Delete(code domain.ObjectCode) error {
	return errors.New("Method not implemented")
}
func (s *mgoObjectStorage) Save(object domain.Object) error {
	return nil
}
