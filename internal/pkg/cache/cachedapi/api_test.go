package cachedapi_test

import (
	"log"

	"github.com/Toggly/core/internal/api"
	in "github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/pkg/cache/cachedapi"
	"github.com/Toggly/core/internal/pkg/engine"
	"github.com/Toggly/core/internal/pkg/storage/mongo"
	"github.com/Toggly/core/pkg/cache"
	"github.com/globalsign/mgo"
)

const MongoTestUrl = "mongodb://localhost:27017/toggly_cache_test"

func getEngineAndCache() (api.TogglyAPI, cache.DataCache) {
	dataStorage, err := mongo.NewMongoStorage(MongoTestUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}
	dataCache := &in.InMemoryCache{
		Storage: make(map[string][]byte, 0),
	}
	if err != nil {
		log.Fatalf(err.Error())
	}
	return cachedapi.NewCachedAPI(engine.NewTogglyAPI(&dataStorage), dataCache), dataCache
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
