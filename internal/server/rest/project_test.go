package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/server/rest"
	asserts "github.com/stretchr/testify/assert"
)

var regDate time.Time

func TestRestProject(t *testing.T) {
	assert := asserts.New(t)
	BeforeTest()

	tt := []struct {
		name      string
		method    string
		path      string
		status    int
		validator func(r *http.Response)
		skip      bool
		body      interface{}
	}{
		{
			name:   "empty list",
			method: http.MethodGet,
			path:   "/project",
			status: http.StatusOK,
			validator: func(r *http.Response) {
				body, err := bodyJSON(r)
				assert.Nil(err)
				assert.Empty(body["projects"])
			},
		},
		{
			name:   "project not found",
			method: http.MethodGet,
			path:   "/project/project1",
			status: http.StatusNotFound,
			validator: func(r *http.Response) {
				body, err := bodyJSON(r)
				assert.Nil(err)
				assert.Equal("Project not found", body["error"])
			},
		},
		{
			name:   "create project bad request",
			method: http.MethodPost,
			path:   "/project",
			status: http.StatusBadRequest,
		},
		{
			name:   "create project",
			method: http.MethodPost,
			path:   "/project",
			status: http.StatusOK,
			body: &rest.ProjectCreateRequest{
				Code:        "project1",
				Description: "Project 1",
				Status:      domain.ProjectStatusActive,
			},
			validator: func(r *http.Response) {
				b := &domain.Project{}
				err := parseBodyTo(r, b)
				assert.Nil(err)
				assert.Equal(domain.ProjectCode("project1"), b.Code)
				assert.Equal("Project 1", b.Description)
				assert.Equal(ow, b.OwnerID)
				assert.NotNil(b.RegDate)
				regDate = b.RegDate
				assert.Equal(domain.ProjectStatusActive, b.Status)
			},
		},
		{
			name:   "projects list",
			method: http.MethodGet,
			path:   "/project",
			status: http.StatusOK,
			validator: func(r *http.Response) {
				body, err := bodyJSON(r)
				assert.Nil(err)
				assert.Len(body["projects"], 1)
			},
		},
		{
			name:   "create project unique index error",
			method: http.MethodPost,
			path:   "/project",
			status: http.StatusBadRequest,
			body: &rest.ProjectCreateRequest{
				Code:        "project1",
				Description: "Project 1",
				Status:      domain.ProjectStatusActive,
			},
			validator: func(r *http.Response) {
				body, err := bodyJSON(r)
				assert.Nil(err)
				assert.Contains(body["error"], "Unique index error:")
			},
		},
		{
			name:   "get project",
			method: http.MethodGet,
			path:   "/project/project1",
			status: http.StatusOK,
			validator: func(r *http.Response) {
				b := &domain.Project{}
				err := parseBodyTo(r, b)
				assert.Nil(err)
				assert.Equal(domain.ProjectCode("project1"), b.Code)
				assert.Equal("Project 1", b.Description)
				assert.Equal(ow, b.OwnerID)
				assert.Equal(regDate, b.RegDate)
				assert.Equal(domain.ProjectStatusActive, b.Status)
			},
		},
		{
			name:   "update project not found",
			method: http.MethodPut,
			path:   "/project",
			status: http.StatusNotFound,
			body: &rest.ProjectCreateRequest{
				Code:        "project2",
				Description: "Project 2",
				Status:      domain.ProjectStatusDisabled,
			},
		},
		{
			name:   "update project",
			method: http.MethodPut,
			path:   "/project",
			status: http.StatusOK,
			body: &rest.ProjectCreateRequest{
				Code:        "project1",
				Description: "Project 1 Updated",
				Status:      domain.ProjectStatusDisabled,
			},
			validator: func(r *http.Response) {
				b := &domain.Project{}
				err := parseBodyTo(r, b)
				assert.Nil(err)
				assert.Equal(domain.ProjectCode("project1"), b.Code)
				assert.Equal("Project 1 Updated", b.Description)
				assert.Equal(ow, b.OwnerID)
				assert.Equal(regDate, b.RegDate)
				assert.Equal(domain.ProjectStatusDisabled, b.Status)
			},
		},
		{
			name:   "delete project not found",
			method: http.MethodDelete,
			path:   "/project/project2",
			status: http.StatusNotFound,
		},
		{
			name:   "delete project",
			method: http.MethodDelete,
			path:   "/project/project1",
			status: http.StatusOK,
		},
	}

	rs := httptest.NewServer(GetRouter().Router())
	defer rs.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}
			var body []byte
			var err error
			if tc.body != nil {
				body, err = json.Marshal(tc.body)
				assert.Nil(err)
			}
			req, err := http.NewRequest(tc.method, rs.URL+v1Path+tc.path, bytes.NewBuffer(body))
			assert.Nil(err)
			req.Header = http.Header{
				rest.XTogglyAuth:      []string{TestAuthToken},
				rest.XTogglyOwnerID:   []string{ow},
				rest.XTogglyRequestID: []string{"request"},
			}
			r, err := rs.Client().Do(req)
			fmt.Printf("\nRESP: %v\n\n", r)
			assert.Nil(err)
			assert.NotNil(r)
			assert.Equal(tc.status, r.StatusCode)
			assert.Equal("Toggly", r.Header[rest.XServiceName][0])
			assert.Equal("test", r.Header[rest.XServiceVersion][0])
			assert.Equal("request", r.Header[rest.XTogglyRequestID][0])
			assert.Contains(r.Header[http.CanonicalHeaderKey("Content-Type")][0], "application/json")
			if tc.validator != nil {
				tc.validator(r)
			}
		})
	}

	AfterTest()
}
