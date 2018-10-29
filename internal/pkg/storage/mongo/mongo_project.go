package mongo

import (
	"fmt"

	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mgProjectStorage struct {
	owner   string
	session *mgo.Session
}

func (s *mgProjectStorage) List() ([]*domain.Project, error) {
	conn := s.session.Copy()
	defer conn.Close()
	items := make([]*domain.Project, 0)
	err := getCollection(conn, "project").Find(bson.M{"owner": s.owner}).All(&items)
	return items, err
}

func (s *mgProjectStorage) Get(code domain.ProjectCode) (project *domain.Project, err error) {
	conn := s.session.Copy()
	defer conn.Close()
	err = getCollection(conn, "project").Find(bson.M{"owner": s.owner, "code": code}).One(&project)
	if err == mgo.ErrNotFound {
		return nil, storage.ErrNotFound
	}
	return project, err
}

func (s *mgProjectStorage) Delete(code domain.ProjectCode) (err error) {
	conn := s.session.Copy()
	defer conn.Close()
	err = getCollection(conn, "project").Remove(bson.M{"owner": s.owner, "code": code})
	if err == mgo.ErrNotFound {
		return storage.ErrNotFound
	}
	// TODO remove environments for this project

	return err
}

func (s *mgProjectStorage) Save(project *domain.Project) error {
	conn := s.session.Copy()
	defer conn.Close()

	collection := getCollection(conn, "project")
	idx := mgo.Index{
		Key:    []string{"owner", "code"},
		Unique: true,
	}
	collection.EnsureIndex(idx)

	err := collection.Insert(project)
	if err != nil {
		if mgo.IsDup(err) {
			return &storage.UniqueIndexError{
				Type: "Project",
				Key:  fmt.Sprintf("owner: %s, code: %s", project.OwnerID, project.Code),
			}
		}
		return err
	}
	return nil
}

//TODO
func (s *mgProjectStorage) Update(project *domain.Project) error {
	conn := s.session.Copy()
	defer conn.Close()

	collection := getCollection(conn, "project")
	idx := mgo.Index{
		Key:    []string{"owner", "code"},
		Unique: true,
	}
	collection.EnsureIndex(idx)

	err := collection.Update(bson.M{"owner": project.OwnerID, "code": project.Code}, project)
	if err != nil {
		return err
	}
	return nil
}

func (s *mgProjectStorage) For(projectCode domain.ProjectCode) storage.ForProject {
	return &mgForProject{
		projectCode: projectCode,
		session:     s.session,
	}
}

type mgForProject struct {
	projectCode domain.ProjectCode
	session     *mgo.Session
	owner       string
}

func (s *mgForProject) Environments() storage.EnvironmentStorage {
	return &mgoEnvStorage{
		projectCode: s.projectCode,
		session:     s.session,
		owner:       s.owner,
	}
}
