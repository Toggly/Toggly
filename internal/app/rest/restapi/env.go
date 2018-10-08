package restapi

import (
	"log"
	"net/http"

	"github.com/Toggly/core/internal/app/rest"
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/go-chi/chi"
)

// EnvironmentAPI servers objects
type EnvironmentAPI struct {
	Cache   cache.DataCache
	Storage storage.DataStorage
}

// Routes returns routes for environments
func (api *EnvironmentAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Get("/", api.cached(api.list))
		g.Get("/{code}", api.cached(api.getEnvironment))
	})
	return router
}

func (api *EnvironmentAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, api.Cache)
}

func (api *EnvironmentAPI) list(w http.ResponseWriter, r *http.Request) {
	proj := domain.ProjectCode(chi.URLParam(r, "project_code"))
	list, err := api.Storage.ForOwner(rest.OwnerFromContext(r)).Projects().For(proj).Environments().List()
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	if list == nil {
		rest.NotFoundResponse(w, r)
		return
	}
	rest.JSONResponse(w, r, list)
}

func (api *EnvironmentAPI) getEnvironment(w http.ResponseWriter, r *http.Request) {
	envID := domain.EnvironmentCode(chi.URLParam(r, "code"))
	proj := domain.ProjectCode(chi.URLParam(r, "project_code"))
	env, err := api.Storage.ForOwner(rest.OwnerFromContext(r)).Projects().For(proj).Environments().Get(envID)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	if env == nil {
		rest.NotFoundResponse(w, r)
		return
	}
	rest.JSONResponse(w, r, env)
}
