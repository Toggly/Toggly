package storage

import (
	"strings"

	"github.com/Toggly/core/app/data"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// NewMongoStorage implements DataStorage interface for MongoDB
func NewMongoStorage(url string) (DataStorage, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	return &mgStorage{
		session: session,
	}, nil
}

type mgStorage struct {
	session *mgo.Session
}

func (s *mgStorage) Projects(ownerID string) ProjectStorage {
	return &mgProjectStorage{
		owner:   ownerID,
		storage: s,
	}
}

type mgProjectStorage struct {
	owner   string
	storage *mgStorage
}

func (s *mgProjectStorage) List() (items []*data.Project, err error) {
	conn := s.storage.session.Copy()
	defer conn.Close()

	err = conn.DB("").C("project").Find(bson.M{"owner": s.owner}).All(&items)
	return items, err
}

func (s *mgProjectStorage) Get(code data.ProjectCode) (*data.Project, error) {
	return nil, nil
}

func (s *mgProjectStorage) Save(project data.Project) error {
	conn := s.storage.session.Copy()
	defer conn.Close()
	project.OwnerID = s.owner
	collection := conn.DB("").C("project")
	idx := mgo.Index{
		Key:    []string{"owner", "code"},
		Unique: true,
	}
	collection.EnsureIndex(idx)

	err := collection.Insert(project)
	if err != nil && strings.Contains(err.Error(), "E11000") {
		return &UniqueIndexError{err.Error()}
	}
	return err
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
