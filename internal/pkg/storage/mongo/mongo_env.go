package mongo

import (
	"fmt"

	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mgoEnvStorage struct {
	projectCode domain.ProjectCode
	session     *mgo.Session
	owner       string
}

func (s *mgoEnvStorage) List() ([]*domain.Environment, error) {
	conn := s.session.Copy()
	defer conn.Close()
	items := make([]*domain.Environment, 0)
	err := getCollection(conn, "env").Find(bson.M{"project_code": s.projectCode}).All(&items)
	return items, err
}

func (s *mgoEnvStorage) Get(code domain.EnvironmentCode) (env *domain.Environment, err error) {
	conn := s.session.Copy()
	defer conn.Close()
	err = getCollection(conn, "env").Find(bson.M{"project_code": s.projectCode, "code": code}).One(&env)
	if err == mgo.ErrNotFound {
		return nil, storage.ErrNotFound
	}
	return env, err
}

func (s *mgoEnvStorage) Delete(code domain.EnvironmentCode) (err error) {
	conn := s.session.Copy()
	defer conn.Close()
	err = getCollection(conn, "env").Remove(bson.M{"project_code": s.projectCode, "code": code})
	if err == mgo.ErrNotFound {
		return storage.ErrNotFound
	}
	return err
}

func ensureEnvIndex(collection *mgo.Collection) {
	idx := mgo.Index{
		Key:    []string{"project_code", "code"},
		Unique: true,
	}
	collection.EnsureIndex(idx)
}

func (s *mgoEnvStorage) Save(env *domain.Environment) error {
	conn := s.session.Copy()
	defer conn.Close()

	collection := getCollection(conn, "env")
	ensureEnvIndex(collection)

	err := collection.Insert(env)
	if err != nil {
		if mgo.IsDup(err) {
			return &storage.UniqueIndexError{
				Type: "Environment",
				Key:  fmt.Sprintf("project_code: %s, code: %s", env.ProjectCode, env.Code),
			}
		}
		return err
	}
	return nil
}

func (s *mgoEnvStorage) Update(env *domain.Environment) error {
	conn := s.session.Copy()
	defer conn.Close()

	collection := getCollection(conn, "env")
	ensureEnvIndex(collection)

	err := collection.Update(bson.M{"project_code": env.ProjectCode, "code": env.Code}, env)
	if err != nil {
		return err
	}
	return nil
}

func (s *mgoEnvStorage) For(code domain.EnvironmentCode) storage.ForEnvironment {
	return &mgoForEnvironment{
		projectCode: s.projectCode,
		env:         code,
		session:     s.session,
		owner:       s.owner,
	}
}

type mgoForEnvironment struct {
	projectCode domain.ProjectCode
	env         domain.EnvironmentCode
	session     *mgo.Session
	owner       string
}

func (s *mgoForEnvironment) Objects() storage.ObjectStorage {
	return &mgoObjectStorage{
		projectCode: s.projectCode,
		envCode:     s.env,
		session:     s.session,
		owner:       s.owner,
	}
}
