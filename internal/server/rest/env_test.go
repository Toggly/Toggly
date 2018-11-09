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
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Empty(b["environments"])
			},
		},
		{
			name:   "List env but project not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project2/env",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
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
				b, err := bodyJSON(body)
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
				b, err := bodyJSON(body)
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
				b, err := bodyJSON(body)
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
				b, err := bodyJSON(body)
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
				b, err := bodyJSON(body)
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
