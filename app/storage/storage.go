package storage

import "github.com/Toggly/backend/app/data"

// DataStorage defines storage interface for dictionary
type DataStorage interface {
	ListProjects(env data.ObjectID) []data.Project
	GetProject(id data.ObjectID) data.Project
	GetObject(project data.ObjectID, env data.ObjectID, obj data.ObjectCode) data.Object
}
