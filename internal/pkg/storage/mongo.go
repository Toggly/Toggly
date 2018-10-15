package storage

import (
	"github.com/globalsign/mgo"
)

// NewMongoStorage implements DataStorage interface for MongoDB
func NewMongoStorage(url string) (DataStorage, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &mgStorage{
		session: session,
	}, nil
}

type mgStorage struct {
	session *mgo.Session
}

func (s *mgStorage) ForOwner(ownerID string) OwnerStorage {
	return &mgOwnerStorage{owner: ownerID, storage: s}
}

type mgOwnerStorage struct {
	owner   string
	storage *mgStorage
}

func (s *mgOwnerStorage) Projects() ProjectStorage {
	return &mgProjectStorage{
		owner:   s.owner,
		storage: s.storage,
	}
}
