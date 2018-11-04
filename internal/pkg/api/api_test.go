package api_test

import (
	"github.com/Toggly/core/internal/domain"
	"github.com/globalsign/mgo"
)

const ProjectCode = domain.ProjectCode("p1")

const MongoTestUrl = "mongodb://localhost:27017/toggly_test"

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
