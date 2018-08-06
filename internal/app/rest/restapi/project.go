package restapi

import (
	"log"
	"net/http"

	"github.com/Toggly/core/internal/app/api"
	"github.com/Toggly/core/internal/app/rest"
	"github.com/Toggly/core/internal/pkg/cache"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
)

// ProjectAPI servers project api namespace
type ProjectAPI struct {
	Cache  cache.DataCache
	Engine *api.Engine
}

// Routes returns routes for project namespace
func (p *ProjectAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Get("/", p.cached(p.list))
		group.Get("/{id}", p.cached(p.getProject))
	})
	return router
}

func (p *ProjectAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, p.Cache)
}

func (p *ProjectAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := p.Engine.Project.List(rest.OwnerFromContext(r))
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	if list == nil {
		rest.NotFoundResponse(w, r)
		return
	}
	render.JSON(w, r, list)
}

func (p *ProjectAPI) getProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	proj, err := p.Engine.Project.Get(rest.OwnerFromContext(r), id)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	if proj == nil {
		rest.NotFoundResponse(w, r)
		return
	}
	render.JSON(w, r, proj)
}
