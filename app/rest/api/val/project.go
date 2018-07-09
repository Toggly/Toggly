package val

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"

	"github.com/Toggly/backend/app/cache"
	"github.com/Toggly/backend/app/rest"
	"github.com/go-chi/chi"
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

func find(id string) (*rest.Project, error) {
	for _, v := range projects() {
		if string(v.ID) == id {
			return &v, nil
		}
	}
	return nil, nil
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
		g.Get("/{id}", p.cached(p.getProject))
	})
	return router
}

func (p *ProjectAPI) cached(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return cache.Cached(fn, p.Cache)
}

// GET ../project
func (p *ProjectAPI) list(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, projects())
}

// GET ../project/{id}
func (p *ProjectAPI) getProject(w http.ResponseWriter, r *http.Request) {
	proj, err := find(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, proj)
}
