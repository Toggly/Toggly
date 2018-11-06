package api_test

import (
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/api"
	"github.com/Toggly/core/internal/pkg/storage/mongo"
	"github.com/globalsign/mgo"
)

const ProjectCode = domain.ProjectCode("p1")

const MongoTestUrl = "mongodb://localhost:27017/toggly_test"

const ow = "test_owner"

func GetApi() *api.ProjectAPI {
	dataStorage, _ := mongo.NewMongoStorage(MongoTestUrl)
	engine := &api.Engine{Storage: &dataStorage}
	pApi := engine.ForOwner(ow).Projects()
	return pApi
}

func DropDB() {
	session, _ := mgo.Dial(MongoTestUrl)
	session.DB("").DropDatabase()
}

func BeforeTest() {
	DropDB()
}

func AfterTest() {
	DropDB()
}
