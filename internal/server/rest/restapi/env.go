package restapi

import (
	"log"
	"net/http"

	"github.com/Toggly/core/internal/pkg/api"
	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/server/rest"
	"github.com/go-chi/chi"
)

// EnvironmentRestAPI servers objects
type EnvironmentRestAPI struct {
	Cache  cache.DataCache
	Engine *api.Engine
}

// Routes returns routes for environments
func (a *EnvironmentRestAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Get("/", a.cached(a.list))
		// g.Get("/{code}", a.cached(a.getEnvironment))
	})
	return router
}

func (a *EnvironmentRestAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, a.Cache)
}

func (a *EnvironmentRestAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := a.Engine.ForOwner(owner(r)).Project().For(projectCode(r)).List()
	if err != nil {
		switch err {
		case api.ErrNotFound:
			rest.NotFoundResponse(w, r)
		default:
			log.Printf("[ERROR] %v", err)
			rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	rest.JSONResponse(w, r, list)
}
