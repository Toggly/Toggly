package mongo

import (
	"fmt"

	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mgEnvironmentStorage struct {
	project domain.ProjectCode
	session *mgo.Session
	owner   string
}

func (s *mgEnvironmentStorage) List() ([]*domain.Environment, error) {
	conn := s.session.Copy()
	defer conn.Close()
	items := make([]*domain.Environment, 0)
	err := getCollection(conn, "env").Find(bson.M{"project_code": s.project}).All(&items)
	return items, err
}

func (s *mgEnvironmentStorage) Get(code domain.EnvironmentCode) (env *domain.Environment, err error) {
	conn := s.session.Copy()
	defer conn.Close()
	err = getCollection(conn, "env").Find(bson.M{"project_code": s.project, "code": code}).One(&env)
	if err == mgo.ErrNotFound {
		return nil, storage.ErrNotFound
	}
	return env, err
}

func (s *mgEnvironmentStorage) Delete(code domain.EnvironmentCode) (err error) {
	conn := s.session.Copy()
	defer conn.Close()
	err = getCollection(conn, "env").Remove(bson.M{"project_code": s.project, "code": code})
	if err == mgo.ErrNotFound {
		return storage.ErrNotFound
	}
	// TODO remove objects for environment
	return err
}

func (s *mgEnvironmentStorage) Save(env *domain.Environment) (*domain.Environment, error) {
	conn := s.session.Copy()
	defer conn.Close()

	collection := getCollection(conn, "env")
	idx := mgo.Index{
		Key:    []string{"project_code", "code"},
		Unique: true,
	}
	collection.EnsureIndex(idx)

	err := collection.Insert(env)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, &storage.UniqueIndexError{
				Type: "Environment",
				Key:  fmt.Sprintf("project_code:%s, code: %s", env.ProjectCode, env.Code),
			}
		}
		return nil, err
	}
	return env, nil
}

func (s *mgEnvironmentStorage) For(code domain.EnvironmentCode) storage.ForEnvironment {
	return &mgForEnvironment{}
}

type mgForEnvironment struct {
	project domain.ProjectCode
	env     domain.EnvironmentCode
	session *mgo.Session
	owner   string
}

func (s *mgForEnvironment) Objects() storage.ObjectStorage {
	return &mgObjectStorage{}
}
