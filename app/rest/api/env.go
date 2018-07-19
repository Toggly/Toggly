package api

import (
	"log"
	"net/http"

	"github.com/Toggly/core/app/data"
	"github.com/Toggly/core/app/rest"
	"github.com/Toggly/core/app/rest/cache"
	"github.com/Toggly/core/app/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// EnvironmentAPI servers objects
type EnvironmentAPI struct {
	Cache   cache.DataCache
	Storage storage.DataStorage
}

// Routes returns routes for environments
func (p *EnvironmentAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Get("/", p.cached(p.list))
		g.Get("/{code}", p.cached(p.getEnvironment))
	})
	return router
}

func (p *EnvironmentAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return cache.Cached(fn, p.Cache)
}

func (p *EnvironmentAPI) list(w http.ResponseWriter, r *http.Request) {
	proj := data.ProjectCode(chi.URLParam(r, "project_code"))
	list, err := p.Storage.Projects(rest.CtxOwner(r)).For(proj).Environments().List()
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

func (p *EnvironmentAPI) getEnvironment(w http.ResponseWriter, r *http.Request) {
	envID := data.EnvironmentCode(chi.URLParam(r, "code"))
	proj := data.ProjectCode(chi.URLParam(r, "project_code"))
	env, err := p.Storage.Projects(rest.CtxOwner(r)).For(proj).Environments().Get(envID)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	if env == nil {
		rest.NotFoundResponse(w, r)
		return
	}
	render.JSON(w, r, env)
}
