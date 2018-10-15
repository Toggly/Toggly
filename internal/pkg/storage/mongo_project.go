package storage

import (
	"fmt"

	"github.com/Toggly/core/internal/domain"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mgProjectStorage struct {
	owner   string
	storage *mgStorage
}

func getCollection(conn *mgo.Session, name string) *mgo.Collection {
	return conn.DB("").C(name)
}

func (s *mgProjectStorage) List() ([]*domain.Project, error) {
	conn := s.storage.session.Copy()
	defer conn.Close()
	items := make([]*domain.Project, 0)
	err := getCollection(conn, "project").Find(bson.M{"owner": s.owner}).All(&items)
	return items, err
}

func (s *mgProjectStorage) Get(code domain.ProjectCode) (project *domain.Project, err error) {
	conn := s.storage.session.Copy()
	defer conn.Close()
	collection := getCollection(conn, "project")
	err = collection.Find(bson.M{"owner": s.owner, "code": code}).One(&project)
	if err == mgo.ErrNotFound {
		return nil, ErrNotFound
	}
	return project, err
}

func (s *mgProjectStorage) Delete(code domain.ProjectCode) (err error) {
	conn := s.storage.session.Copy()
	defer conn.Close()
	collection := getCollection(conn, "project")
	err = collection.Remove(bson.M{"owner": s.owner, "code": code})
	if err == mgo.ErrNotFound {
		return ErrNotFound
	}
	// TODO remove environments for this project

	return err
}

func (s *mgProjectStorage) Save(project *domain.Project) (*domain.Project, error) {
	conn := s.storage.session.Copy()
	defer conn.Close()

	proj := &domain.Project{
		OwnerID:     s.owner,
		Code:        project.Code,
		Description: project.Description,
		RegDate:     project.RegDate,
		Status:      project.Status,
	}

	collection := getCollection(conn, "project")
	idx := mgo.Index{
		Key:    []string{"owner", "code"},
		Unique: true,
	}
	collection.EnsureIndex(idx)

	err := collection.Insert(proj)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, &UniqueIndexError{
				Type: "Project",
				Key:  fmt.Sprintf("owner:%s, code: %s", proj.OwnerID, proj.Code),
			}
		}
		return nil, err
	}
	return proj, nil
}

func (s *mgProjectStorage) For(project domain.ProjectCode) ForProject {
	return &mgForProject{
		project: project,
	}
}
