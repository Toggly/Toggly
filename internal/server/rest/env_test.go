package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/server/rest"
	asserts "github.com/stretchr/testify/assert"
)

var regDateEnv time.Time

func TestRestEnvironment(t *testing.T) {
	assert := asserts.New(t)
	BeforeTest()

	tt := []TestCase{
		{
			name:   "empty list",
			method: http.MethodGet,
			path:   "/api/v1/project/project1/env",
			status: http.StatusOK,
			validator: func(body []byte) {
				var b []*domain.Environment
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Empty(b)
			},
		},
		{
			name:   "List env but project not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project2/env",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Equal("Project not found", b["error"])
			},
		},
		{
			name:   "Env not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project1/env/env1",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Equal("Environment not found", b["error"])
			},
		},
		{
			name:   "Get env but project not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project2/env/env1",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Equal("Project not found", b["error"])
			},
		},
		{
			name:   "Create env but project not found",
			method: http.MethodPost,
			path:   "/api/v1/project/project2/env",
			body:   &rest.EnvironmentCreateRequest{},
			status: http.StatusNotFound,
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Equal("Project not found", b["error"])
			},
		},
		{
			name:   "Update env but project not found",
			method: http.MethodPut,
			path:   "/api/v1/project/project2/env",
			body:   &rest.EnvironmentCreateRequest{},
			status: http.StatusNotFound,
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Equal("Project not found", b["error"])
			},
		},
		{
			name:   "Delete env but project not found",
			method: http.MethodDelete,
			path:   "/api/v1/project/project2/env/env1",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Equal("Project not found", b["error"])
			},
		},
		{
			name:   "Create env bad request",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env",
			status: http.StatusBadRequest,
		},
		{
			name:   "Create env bad request",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env",
			body:   &rest.EnvironmentCreateRequest{},
			status: http.StatusBadRequest,
		},
		{
			name:   "Create env",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env",
			body: &rest.EnvironmentCreateRequest{
				Code:        "env1",
				Description: "Env description",
				Protected:   false,
			},
			status: http.StatusOK,
			validator: func(body []byte) {
				env := &domain.Environment{}
				err := parseBodyTo(body, env)
				assert.Equal(domain.EnvironmentCode("env1"), env.Code)
				assert.Equal("Env description", env.Description)
				assert.Equal(ow, env.OwnerID)
				assert.Equal(domain.ProjectCode("project1"), env.ProjectCode)
				assert.Equal(false, env.Protected)
				assert.NotNil(env.RegDate)
				regDateEnv = env.RegDate
				assert.Nil(err)
			},
		},
		{
			name:   "Create env unique index error",
			method: http.MethodPost,
			path:   "/api/v1/project/project1/env",
			body: &rest.EnvironmentCreateRequest{
				Code:        "env1",
				Description: "Env description",
				Protected:   false,
			},
			status: http.StatusBadRequest,
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Contains(b["error"], "Unique index error:")
			},
		},
		{
			name:   "Update env bad request",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env",
			status: http.StatusBadRequest,
		},
		{
			name:   "Update env bad request",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env",
			body:   &rest.EnvironmentCreateRequest{},
			status: http.StatusBadRequest,
		},
		{
			name:   "Update env but env not found",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env",
			body: &rest.EnvironmentCreateRequest{
				Code:        "env2",
				Description: "Env description 2",
				Protected:   false,
			},
			status: http.StatusNotFound,
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Equal("Environment not found", b["error"])
			},
		},
		{
			name:   "Update env",
			method: http.MethodPut,
			path:   "/api/v1/project/project1/env",
			body: &rest.EnvironmentCreateRequest{
				Code:        "env1",
				Description: "Env description 1",
				Protected:   true,
			},
			status: http.StatusOK,
			validator: func(body []byte) {
				env := &domain.Environment{}
				err := parseBodyTo(body, env)
				assert.Equal(domain.EnvironmentCode("env1"), env.Code)
				assert.Equal("Env description 1", env.Description)
				assert.Equal(ow, env.OwnerID)
				assert.Equal(domain.ProjectCode("project1"), env.ProjectCode)
				assert.Equal(true, env.Protected)
				assert.Equal(regDateEnv, env.RegDate)
				assert.Nil(err)
			},
		},
		{
			name:   "Delete env but env not found",
			method: http.MethodDelete,
			path:   "/api/v1/project/project1/env/env2",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Equal("Environment not found", b["error"])
			},
		},
		{
			name:   "delete not empty environment",
			method: http.MethodDelete,
			path:   "/api/v1/project/project1/env/env1",
			status: http.StatusLocked,
			before: func(rs *httptest.Server) {
				body, err := json.Marshal(&rest.ObjectCreateRequest{Code: "obj1"})
				assert.Nil(err)
				req, err := http.NewRequest(http.MethodPost, rs.URL+"/api/v1/project/project1/env/env1/object", bytes.NewBuffer(body))
				assert.Nil(err)
				req.Header = http.Header{
					rest.XTogglyAuth:    []string{TestAuthToken},
					rest.XTogglyOwnerID: []string{ow},
				}
				rs.Client().Do(req)
			},
			validator: func(body []byte) {
				var b map[string]interface{}
				err := parseBodyTo(body, &b)
				assert.Nil(err)
				assert.Equal(rest.ErrEnvironmentNotEmpty, b["error"])
			},
			after: func(rs *httptest.Server) {
				req, err := http.NewRequest(http.MethodDelete, rs.URL+"/api/v1/project/project1/env/env1/object/obj1", nil)
				assert.Nil(err)
				req.Header = http.Header{
					rest.XTogglyAuth:    []string{TestAuthToken},
					rest.XTogglyOwnerID: []string{ow},
				}
				rs.Client().Do(req)
			},
		}, {
			name:   "Delete env",
			method: http.MethodDelete,
			path:   "/api/v1/project/project1/env/env1",
			status: http.StatusOK,
		},
	}

	rs := httptest.NewServer(GetRouter().Router())
	defer rs.Close()

	body, err := json.Marshal(&rest.ProjectCreateRequest{
		Code:        "project1",
		Description: "Project 1",
		Status:      domain.ProjectStatusActive,
	})
	assert.Nil(err)
	req, err := http.NewRequest(http.MethodPost, rs.URL+"/api/v1/project", bytes.NewBuffer(body))
	assert.Nil(err)
	req.Header = http.Header{
		rest.XTogglyAuth:      []string{TestAuthToken},
		rest.XTogglyOwnerID:   []string{ow},
		rest.XTogglyRequestID: []string{"request"},
	}
	_, err = rs.Client().Do(req)
	assert.Nil(err)

	for _, tc := range tt {
		runTestCase(t, rs, tc)
	}

	AfterTest()

}
