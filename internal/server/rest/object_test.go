package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/server/rest"
	asserts "github.com/stretchr/testify/assert"
)

func TestRestObject(t *testing.T) {
	assert := asserts.New(t)
	BeforeTest()

	tt := []TestCase{
		// Project/Env not found LIST
		{
			name:   "List objects, but project not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project2/env/env1/object",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Project not found", b["error"])
			},
		},
		{
			name:   "List objects, but env not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project1/env/env2/object",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Environment not found", b["error"])
			},
		},
		// Project/Env not found GET
		{
			name:   "Get object, but project not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project2/env/env1/object/obj1",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Project not found", b["error"])
			},
		},
		{
			name:   "Get object, but env not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project1/env/env2/object/obj1",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Environment not found", b["error"])
			},
		},
		// Project/Env not found CREATE
		{
			name:   "Get object, but project not found",
			method: http.MethodPost,
			path:   "/api/v1/project/project2/env/env1/object",
			body: &rest.ObjectCreateRequest{
				Code:        "obj1",
				Description: "Obj1",
			},
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Project not found", b["error"])
			},
		},
		{
			name:   "Get object, but env not found",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env/env2/object",
			body: &rest.ObjectCreateRequest{
				Code:        "obj1",
				Description: "Obj1",
			},
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Environment not found", b["error"])
			},
		},
		// Project/Env not found DELETE
		{
			name:   "Delete object, but project not found",
			method: http.MethodDelete,
			path:   "/api/v1/project/project2/env/env1/object/obj1",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Project not found", b["error"])
			},
		},
		{
			name:   "Delete object, but env not found",
			method: http.MethodDelete,
			path:   "/api/v1/project/project1/env/env2/object/obj1",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Environment not found", b["error"])
			},
		},
		// Object tests
		{
			name:   "empty list",
			method: http.MethodGet,
			path:   "/api/v1/project/project1/env/env1/object",
			status: http.StatusOK,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Empty(b["objects"])
			},
		},
		{
			name:   "Get object not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project1/env/env1/object/obj1",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Object not found", b["error"])
			},
		},
		{
			name:   "Create object bad request",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env/env1/object",
			status: http.StatusBadRequest,
		},
		{
			name:   "Create object bad request",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env/env1/object",
			body:   &rest.ObjectCreateRequest{},
			status: http.StatusBadRequest,
		},
		{
			name:   "Create object bad request no parrent",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env/env1/object",
			body: &rest.ObjectCreateRequest{
				Code: "obj1",
				Inherits: &domain.ObjectInheritance{
					ProjectCode: "project1",
					EnvCode:     "env1",
					ObjectCode:  "none_code",
				},
			},
			status: http.StatusBadRequest,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("object parrent does not exists", b["error"])
			},
		},
		{
			name:   "Create object",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env/env1/object",
			body: &rest.ObjectCreateRequest{
				Code:        "obj1",
				Description: "Object 1",
				Parameters: []*domain.Parameter{
					&domain.Parameter{
						Code:        "param1",
						Description: "Par Desc 1",
						Type:        domain.ParameterBool,
						Value:       true,
					},
				},
			},
			status: http.StatusOK,
			validator: func(body []byte) {
				obj := &domain.Object{}
				err := parseBodyTo(body, obj)
				assert.Nil(err)
			},
		},
		{
			name:   "Create object unique index error",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env/env1/object",
			body: &rest.ObjectCreateRequest{
				Code: "obj1",
			},
			status: http.StatusBadRequest,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Unique index error: Object [env_code: env1, code: obj1]", b["error"])
			},
		},
		{
			name:   "Create object type mismatch",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env/env1/object",
			body: &rest.ObjectCreateRequest{
				Code:        "obj2",
				Description: "Object 2",
				Inherits: &domain.ObjectInheritance{
					ProjectCode: "project1",
					EnvCode:     "env1",
					ObjectCode:  "obj1",
				},
				Parameters: []*domain.Parameter{
					&domain.Parameter{
						Code:        "param1",
						Description: "Par Desc 1",
						Type:        domain.ParameterString,
						Value:       "value",
					},
				},
			},
			status: http.StatusBadRequest,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("object inheritor type mismatch", b["error"])
			},
		},
		{
			name:   "Create inheritor",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env/env1/object",
			body: &rest.ObjectCreateRequest{
				Code:        "obj2",
				Description: "Object 2",
				Inherits: &domain.ObjectInheritance{
					ProjectCode: "project1",
					EnvCode:     "env1",
					ObjectCode:  "obj1",
				},
				Parameters: []*domain.Parameter{
					&domain.Parameter{
						Code:        "param2",
						Description: "Par Desc 2",
						Type:        domain.ParameterString,
						Value:       "value",
					},
				},
			},
			status: http.StatusOK,
		},
		{
			name:   "Update object bad request",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env/env1/object",
			status: http.StatusBadRequest,
		},
		{
			name:   "Update object bad request",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env/env1/object",
			body:   &rest.ObjectCreateRequest{},
			status: http.StatusBadRequest,
		},
		{
			name:   "Update object bad request: parameter type changed",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env/env1/object",
			body: &rest.ObjectCreateRequest{
				Code: "obj1",
				Parameters: []*domain.Parameter{
					&domain.Parameter{
						Code:        "param1",
						Description: "Par Desc 1",
						Type:        domain.ParameterString,
						Value:       "value",
					},
				},
			},
			status: http.StatusBadRequest,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Object parameter `param1` error: Object parameter type changing restricted", b["error"])
			},
		},
		{
			name:   "Update object bad request: Object parameter exists in inheritor",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env/env1/object",
			body: &rest.ObjectCreateRequest{
				Code: "obj1",
				Parameters: []*domain.Parameter{
					&domain.Parameter{
						Code:        "param2",
						Description: "Par Desc 2",
						Type:        domain.ParameterString,
						Value:       "value",
					},
				},
			},
			status: http.StatusBadRequest,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Object parameter `param2` error: Object parameter exists in inheritor: project1:env1:obj2", b["error"])
			},
		},
		{
			name:   "Update object not found",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env/env1/object",
			body: &rest.ObjectCreateRequest{
				Code: "obj123",
			},
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Object not found", b["error"])
			},
		},
		{
			name:   "Delete object not found",
			method: http.MethodDelete,
			path:   "/api/v1/project/project1/env/env1/object/obj999",
			status: http.StatusNotFound,
		},
		{
			name:   "Delete object fail: has inheritor",
			method: http.MethodDelete,
			path:   "/api/v1/project/project1/env/env1/object/obj1",
			status: http.StatusLocked,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal("Object has inheritors", b["error"])
			},
		},
		// {
		// 	name:   "Delete object",
		// 	method: http.MethodDelete,
		// 	path:   "/api/v1/project/project1/env/env1/object/obj1",
		// 	status: http.StatusOK,
		// },
	}

	rs := httptest.NewServer(GetRouter().Router())
	defer rs.Close()

	headers := http.Header{
		rest.XTogglyAuth:      []string{TestAuthToken},
		rest.XTogglyOwnerID:   []string{ow},
		rest.XTogglyRequestID: []string{"request"},
	}

	body, err := json.Marshal(&rest.ProjectCreateRequest{
		Code:        "project1",
		Description: "Project 1",
		Status:      domain.ProjectStatusActive,
	})
	assert.Nil(err)

	req, err := http.NewRequest(http.MethodPost, rs.URL+"/api/v1/project", bytes.NewBuffer(body))
	assert.Nil(err)
	req.Header = headers
	_, err = rs.Client().Do(req)
	assert.Nil(err)

	body, err = json.Marshal(&rest.EnvironmentCreateRequest{
		Code:        "env1",
		Description: "Environment 1",
		Protected:   false,
	})
	assert.Nil(err)

	req, err = http.NewRequest(http.MethodPost, rs.URL+"/api/v1/project/project1/env", bytes.NewBuffer(body))
	assert.Nil(err)
	req.Header = headers
	_, err = rs.Client().Do(req)
	assert.Nil(err)

	for _, tc := range tt {
		runTestCase(t, rs, tc)
	}

	// AfterTest()

}
