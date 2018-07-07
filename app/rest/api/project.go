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

// Routes returns routes for project namespace
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
	key := r.URL.String()
	data, err := p.Cache.Get(key, func() (interface{}, error) {
		id := chi.URLParam(r, "id")
		var p *rest.Project
		for _, v := range projects() {
			if string(v.ID) == id {
				p = &v
				return p, nil
			}
		}
		return nil, nil
	})
	if err != nil {
		render.Status(r, 500)
		render.PlainText(w, r, err.Error())
	}
	if data == nil {
		render.Status(r, 404)
	}
	render.JSON(w, r, data)
}
