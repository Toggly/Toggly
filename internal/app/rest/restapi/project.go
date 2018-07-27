package restapi

import (
	"log"
	"net/http"

	"github.com/Toggly/core/internal/app/rest"
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/pkg/storage"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
)

// ProjectAPI servers project api namespace
type ProjectAPI struct {
	Cache   cache.DataCache
	Storage storage.DataStorage
}

// Routes returns routes for project namespace
func (p *ProjectAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		// g.Get("/", p.list)
		// g.Get("/{id}", p.getProject)
		g.Get("/", p.cached(p.list))
		// g.Get("/{id}", p.cached(p.getProject))
	})
	return router
}

func (p *ProjectAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, p.Cache)
}

func (p *ProjectAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := p.Storage.Projects(rest.CtxOwner(r)).List()
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
	id := domain.ProjectCode(chi.URLParam(r, "id"))
	proj, err := p.Storage.Projects(rest.CtxOwner(r)).Get(id)
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
