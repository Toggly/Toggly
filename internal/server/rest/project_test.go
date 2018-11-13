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

var regDateProject time.Time

func TestRestProject(t *testing.T) {
	assert := asserts.New(t)
	BeforeTest()

	tt := []TestCase{
		{
			name:   "empty list",
			method: http.MethodGet,
			path:   "/api/v1/project",
			status: http.StatusOK,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Empty(b["projects"])
			},
		},
		{
			name:   "project not found",
			method: http.MethodGet,
			path:   "/api/v1/project/project1",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal(rest.ErrProjectNotFound, b["error"])
			},
		},
		{
			name:   "create project bad request",
			method: http.MethodPost,
			path:   "/api/v1/project",
			status: http.StatusBadRequest,
		},
		{
			name:   "create project bad request",
			method: http.MethodPost,
			path:   "/api/v1/project",
			body:   &rest.ProjectCreateRequest{},
			status: http.StatusBadRequest,
		},
		{
			name:   "create project",
			method: http.MethodPost,
			path:   "/api/v1/project",
			status: http.StatusOK,
			body: &rest.ProjectCreateRequest{
				Code:        "project1",
				Description: "Project 1",
				Status:      domain.ProjectStatusActive,
			},
			validator: func(body []byte) {
				b := &domain.Project{}
				err := parseBodyTo(body, b)
				assert.Nil(err)
				assert.Equal(domain.ProjectCode("project1"), b.Code)
				assert.Equal("Project 1", b.Description)
				assert.Equal(ow, b.OwnerID)
				assert.NotNil(b.RegDate)
				regDateProject = b.RegDate
				assert.Equal(domain.ProjectStatusActive, b.Status)
			},
		},
		{
			name:   "projects list",
			method: http.MethodGet,
			path:   "/api/v1/project",
			status: http.StatusOK,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Len(b["projects"], 1)
			},
		},
		{
			name:   "create project unique index error",
			method: http.MethodPost,
			path:   "/api/v1/project",
			status: http.StatusBadRequest,
			body: &rest.ProjectCreateRequest{
				Code:        "project1",
				Description: "Project 1",
				Status:      domain.ProjectStatusActive,
			},
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Contains(b["error"], "Unique index error:")
			},
		},
		{
			name:   "get project",
			method: http.MethodGet,
			path:   "/api/v1/project/project1",
			status: http.StatusOK,
			validator: func(body []byte) {
				b := &domain.Project{}
				err := parseBodyTo(body, b)
				assert.Nil(err)
				assert.Equal(domain.ProjectCode("project1"), b.Code)
				assert.Equal("Project 1", b.Description)
				assert.Equal(ow, b.OwnerID)
				assert.Equal(regDateProject, b.RegDate)
				assert.Equal(domain.ProjectStatusActive, b.Status)
			},
		},
		{
			name:   "update project not found",
			method: http.MethodPut,
			path:   "/api/v1/project",
			status: http.StatusNotFound,
			body: &rest.ProjectCreateRequest{
				Code:        "project2",
				Description: "Project 2",
				Status:      domain.ProjectStatusDisabled,
			},
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal(rest.ErrProjectNotFound, b["error"])
			},
		},
		{
			name:   "update project",
			method: http.MethodPut,
			path:   "/api/v1/project",
			status: http.StatusOK,
			body: &rest.ProjectCreateRequest{
				Code:        "project1",
				Description: "Project 1 Updated",
				Status:      domain.ProjectStatusDisabled,
			},
			validator: func(body []byte) {
				b := &domain.Project{}
				err := parseBodyTo(body, b)
				assert.Nil(err)
				assert.Equal(domain.ProjectCode("project1"), b.Code)
				assert.Equal("Project 1 Updated", b.Description)
				assert.Equal(ow, b.OwnerID)
				assert.Equal(regDateProject, b.RegDate)
				assert.Equal(domain.ProjectStatusDisabled, b.Status)
			},
		},
		{
			name:   "delete project not found",
			method: http.MethodDelete,
			path:   "/api/v1/project/project2",
			status: http.StatusNotFound,
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal(rest.ErrProjectNotFound, b["error"])
			},
		},
		{
			name:   "delete not empty project",
			method: http.MethodDelete,
			path:   "/api/v1/project/project1",
			status: http.StatusLocked,
			before: func(rs *httptest.Server) {
				body, err := json.Marshal(&rest.EnvironmentCreateRequest{
					Code:        "env1",
					Description: "Env description",
					Protected:   false,
				})
				assert.Nil(err)
				req, err := http.NewRequest(http.MethodPost, rs.URL+"/api/v1/project/project1/env", bytes.NewBuffer(body))
				assert.Nil(err)
				req.Header = http.Header{
					rest.XTogglyAuth:    []string{TestAuthToken},
					rest.XTogglyOwnerID: []string{ow},
				}
				rs.Client().Do(req)
			},
			validator: func(body []byte) {
				b, err := bodyJSON(body)
				assert.Nil(err)
				assert.Equal(rest.ErrProjectNotEmpty, b["error"])
			},
			after: func(rs *httptest.Server) {
				req, err := http.NewRequest(http.MethodDelete, rs.URL+"/api/v1/project/project1/env/env1", nil)
				assert.Nil(err)
				req.Header = http.Header{
					rest.XTogglyAuth:    []string{TestAuthToken},
					rest.XTogglyOwnerID: []string{ow},
				}
				rs.Client().Do(req)
			},
		},
		{
			name:   "delete project",
			method: http.MethodDelete,
			path:   "/api/v1/project/project1",
			status: http.StatusOK,
		},
	}

	rs := httptest.NewServer(GetRouter().Router())
	defer rs.Close()

	for _, tc := range tt {
		runTestCase(t, rs, tc)
	}

	AfterTest()
}
