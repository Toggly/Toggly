package api

import (
	"net/http"
	"time"

	"github.com/Toggly/backend/app/cache"
	"github.com/Toggly/backend/app/rest"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func projects() []rest.Project {
	p := make([]rest.Project, 3)

	p[0] = rest.Project{
		ID:      "1",
		Name:    "Simple Project 1",
		RegDate: time.Now(),
		Status:  0,
	}
	p[1] = rest.Project{
		ID:      "2",
		Name:    "Project 2",
		RegDate: time.Now(),
		Status:  0,
	}
	p[2] = rest.Project{
		ID:      "3",
		Name:    "My Project 3",
		RegDate: time.Now(),
		Status:  1,
	}
	return p
}

// ProjectAPI servers project api namespace
type ProjectAPI struct {
	Cache cache.DataCache
}

// Routes returns the project namespace router
func (p *ProjectAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Get("/", p.list)
		g.Get("/{id}", p.getProject)
	})
	return router
}

func (p *ProjectAPI) list(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, projects())
}

func (p *ProjectAPI) getProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	for _, v := range projects() {
		if string(v.ID) == id {
			render.JSON(w, r, v)
			return
		}
	}
	render.Status(r, 404)
	render.PlainText(w, r, "Not found")
}
