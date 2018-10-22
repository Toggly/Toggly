package restapi

import (
	"log"
	"net/http"

	"github.com/Toggly/core/internal/pkg/api"
	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/server/rest"
	"github.com/go-chi/chi"
)

// EnvironmentAPI servers objects
type EnvironmentAPI struct {
	Cache  cache.DataCache
	Engine *api.Engine
}

// Routes returns routes for environments
func (a *EnvironmentAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Get("/", a.cached(a.list))
		// g.Get("/{code}", a.cached(a.getEnvironment))
	})
	return router
}

func (a *EnvironmentAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, a.Cache)
}

func (a *EnvironmentAPI) list(w http.ResponseWriter, r *http.Request) {
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
