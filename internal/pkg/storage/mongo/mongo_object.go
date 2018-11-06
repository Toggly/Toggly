package mongo

import (
	"fmt"

	"github.com/Toggly/core/internal/pkg/storage"

	"github.com/Toggly/core/internal/domain"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mgoObjectStorage struct {
	projectCode domain.ProjectCode
	envCode     domain.EnvironmentCode
	session     *mgo.Session
	owner       string
}

func (s *mgoObjectStorage) List() ([]*domain.Object, error) {
	conn := s.session.Copy()
	defer conn.Close()
	items := make([]*domain.Object, 0)
	err := getCollection(conn, "object").Find(bson.M{"env_code": s.envCode}).All(&items)
	return items, err
}

func (s *mgoObjectStorage) Get(code domain.ObjectCode) (obj *domain.Object, err error) {
	conn := s.session.Copy()
	defer conn.Close()
	err = getCollection(conn, "object").Find(bson.M{"env_code": s.envCode, "code": code}).One(&obj)
	if err == mgo.ErrNotFound {
		return nil, storage.ErrNotFound
	}
	return obj, nil
}

func (s *mgoObjectStorage) ListInheritors(code domain.ObjectCode) ([]*domain.Object, error) {
	conn := s.session.Copy()
	defer conn.Close()

	items := make([]*domain.Object, 0)

	obj, err := s.Get(code)
	if err != nil {
		switch err {
		case storage.ErrNotFound:
			return items, nil
		default:
			return nil, err
		}
	}
	query := bson.M{
		"inherits.project_code": obj.ProjectCode,
		"inherits.env_code":     obj.EnvCode,
		"inherits.object_code":  obj.Code,
	}
	err = getCollection(conn, "object").Find(query).All(&items)
	return items, err
}

func (s *mgoObjectStorage) Delete(code domain.ObjectCode) (err error) {
	conn := s.session.Copy()
	defer conn.Close()
	err = getCollection(conn, "object").Remove(bson.M{"env_code": s.envCode, "code": code})
	if err == mgo.ErrNotFound {
		return storage.ErrNotFound
	}
	return err
}

func ensureObjIndex(collection *mgo.Collection) {
	collection.EnsureIndex(mgo.Index{
		Key:    []string{"env_code", "code"},
		Unique: true,
	})
	collection.EnsureIndex(mgo.Index{
		Key: []string{"inherits.project_code", "inherits.env_code", "inherits.object_code"},
	})
}

func (s *mgoObjectStorage) Save(obj *domain.Object) error {
	conn := s.session.Copy()
	defer conn.Close()

	collection := getCollection(conn, "object")
	ensureObjIndex(collection)

	err := collection.Insert(obj)
	if err != nil {
		if mgo.IsDup(err) {
			return &storage.UniqueIndexError{
				Type: "Object",
				Key:  fmt.Sprintf("env_code: %s, code: %s", obj.EnvCode, obj.Code),
			}
		}
		return err
	}
	return nil
}

func (s *mgoObjectStorage) Update(obj *domain.Object) error {
	conn := s.session.Copy()
	defer conn.Close()

	collection := getCollection(conn, "object")
	ensureObjIndex(collection)

	err := collection.Update(bson.M{"env_code": obj.EnvCode, "code": obj.Code}, obj)
	if err != nil {
		return err
	}
	return nil
}
