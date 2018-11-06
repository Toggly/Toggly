package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Toggly/core/internal/pkg/api"
	"github.com/Toggly/core/internal/pkg/storage/mongo"
	"github.com/Toggly/core/internal/server/rest"
	"github.com/globalsign/mgo"
	asserts "github.com/stretchr/testify/assert"
)

const MongoTestUrl = "mongodb://localhost:27017/toggly_test_rest"
const TestAuthToken = "TestToken"
const ow = "test_owner"
const TestPort = 8081
const v1Path = "/api/v1"

func GetRouter() *rest.APIRouter {
	dataStorage, _ := mongo.NewMongoStorage(MongoTestUrl)
	engine := &api.Engine{Storage: &dataStorage}
	return &rest.APIRouter{
		Version:   "test",
		Cache:     nil,
		Engine:    engine,
		BasePath:  "/api",
		Port:      TestPort,
		IsDebug:   false,
		AuthToken: TestAuthToken,
	}
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

func bodyJSON(r *http.Response) (map[string]interface{}, error) {
	defer r.Body.Close()
	byt, err := ioutil.ReadAll(r.Body)
	fmt.Printf("RAW: %s", byt)
	if err != nil {
		return nil, err
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		return nil, err
	}
	return dat, nil
}

func parseBodyTo(r *http.Response, obj interface{}) error {
	defer r.Body.Close()
	byt, err := ioutil.ReadAll(r.Body)
	fmt.Printf("RAW: %s", byt)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(byt, obj); err != nil {
		return err
	}
	return nil
}

func bodyText(r *http.Response) (string, error) {
	defer r.Body.Close()
	byt, err := ioutil.ReadAll(r.Body)
	fmt.Printf("RAW: %s", byt)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(byt)), nil
}

func TestRestRequestHeaders(t *testing.T) {
	assert := asserts.New(t)
	BeforeTest()

	tt := []struct {
		name        string
		header      http.Header
		status      int
		cType       string
		validator   func(r *http.Response)
		skip        bool
		requestId   string
		resultReqId string
	}{
		{
			name:        "no auth",
			header:      nil,
			status:      http.StatusUnauthorized,
			cType:       "text/plain",
			resultReqId: "^req-\\d*$",
		},
		{
			name:      "wrong token",
			header:    http.Header{rest.XTogglyAuth: []string{"wrong_token"}},
			status:    http.StatusUnauthorized,
			cType:     "text/plain",
			requestId: "12345",
		},
		{
			name:      "authorized but owner not found",
			header:    http.Header{rest.XTogglyAuth: []string{TestAuthToken}},
			status:    http.StatusNotFound,
			requestId: "12345",
			validator: func(r *http.Response) {
				body, err := bodyJSON(r)
				assert.Nil(err)
				assert.Equal("Owner not found", body["error"])
			},
		},
		{
			name:      "authorized",
			header:    http.Header{rest.XTogglyAuth: []string{TestAuthToken}, rest.XTogglyOwnerID: []string{ow}},
			status:    http.StatusNotFound,
			requestId: "12345",
			cType:     "text/plain",
			validator: func(r *http.Response) {
				body, err := bodyText(r)
				assert.Nil(err)
				assert.Equal("404 page not found", body)
			},
		},
	}

	rs := httptest.NewServer(GetRouter().Router())
	defer rs.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}

			req, err := http.NewRequest(http.MethodGet, rs.URL+v1Path, nil)
			assert.Nil(err)
			req.Header = tc.header
			if tc.requestId != "" {
				req.Header[rest.XTogglyRequestID] = []string{tc.requestId}
			}
			r, err := rs.Client().Do(req)
			fmt.Printf("\nRESP: %v\n\n", r)
			assert.Nil(err)
			assert.NotNil(r)
			assert.Equal(tc.status, r.StatusCode)
			assert.Equal("Toggly", r.Header[rest.XServiceName][0])
			assert.Equal("test", r.Header[rest.XServiceVersion][0])
			if tc.requestId != "" {
				assert.Equal(tc.requestId, r.Header[rest.XTogglyRequestID][0])
			} else {
				assert.Regexp(tc.resultReqId, r.Header[rest.XTogglyRequestID][0])
			}
			if tc.cType == "" {
				tc.cType = "application/json"
			}
			assert.Contains(r.Header[http.CanonicalHeaderKey("Content-Type")][0], tc.cType)
			if tc.validator != nil {
				tc.validator(r)
			}
		})
	}

	AfterTest()
}
