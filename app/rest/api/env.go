package api

import (
	"log"
	"net/http"

	"github.com/Toggly/backend/app/data"

	"github.com/Toggly/backend/app/cache"
	"github.com/Toggly/backend/app/rest"
	"github.com/Toggly/backend/app/storage"
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

func (p *EnvironmentAPI) cached(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return cache.Cached(fn, p.Cache)
}

func (p *EnvironmentAPI) list(w http.ResponseWriter, r *http.Request) {
	proj := data.CodeType(chi.URLParam(r, "project_code"))
	list, err := p.Storage.ListEnvironments(proj)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, list)
}

func (p *EnvironmentAPI) getEnvironment(w http.ResponseWriter, r *http.Request) {
	envID := data.CodeType(chi.URLParam(r, "code"))
	proj := data.CodeType(chi.URLParam(r, "project_code"))
	env, err := p.Storage.GetEnvironment(proj, envID)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, env)
}
