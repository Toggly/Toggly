package mongo

import (
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
)

// NewMongoStorage implements DataStorage interface for MongoDB
func NewMongoStorage(url string) (storage.DataStorage, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't connect to %s", url)
	}
	return &mgStorage{
		session: session,
	}, nil
}

func getCollection(conn *mgo.Session, name string) *mgo.Collection {
	return conn.DB("").C(name)
}

type mgStorage struct {
	session *mgo.Session
}

func (s *mgStorage) ForOwner(ownerID string) storage.OwnerStorage {
	return &mgOwnerStorage{owner: ownerID, session: s.session}
}

type mgOwnerStorage struct {
	owner   string
	session *mgo.Session
}

func (s *mgOwnerStorage) Projects() storage.ProjectStorage {
	return &mgProjectStorage{
		owner:   s.owner,
		session: s.session,
	}
}
