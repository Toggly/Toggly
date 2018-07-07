package api

import (
	"net/http"
	"time"

	"github.com/Toggly/backend/app/rest"
	"github.com/go-chi/render"
)

func (a *TogglyAPI) getProject(w http.ResponseWriter, r *http.Request) {
	p := &rest.Project{
		ID:      "1234556",
		Name:    "Simple Project",
		RegDate: time.Now(),
		Status:  0,
	}
	render.JSON(w, r, p)
}
