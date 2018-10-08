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

// ObjectAPI servers objects
type ObjectAPI struct {
	Cache   cache.DataCache
	Storage storage.DataStorage
}

// Routes returns routes for objects
func (api *ObjectAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Get("/", api.cached(api.list))
		g.Get("/{code}", api.cached(api.getObject))
	})
	return router
}

func (api *ObjectAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, api.Cache)
}

func (api *ObjectAPI) list(w http.ResponseWriter, r *http.Request) {
	proj := domain.ProjectCode(chi.URLParam(r, "project_code"))
	env := domain.EnvironmentCode(chi.URLParam(r, "env_code"))
	list, err := api.Storage.ForOwner(rest.OwnerFromContext(r)).Projects().For(proj).Environments().For(env).Objects().List()
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

func (api *ObjectAPI) getObject(w http.ResponseWriter, r *http.Request) {
	proj := domain.ProjectCode(chi.URLParam(r, "project_code"))
	obj := domain.ObjectCode(chi.URLParam(r, "code"))
	env := domain.EnvironmentCode(chi.URLParam(r, "env_code"))
	o, err := api.Storage.ForOwner(rest.OwnerFromContext(r)).Projects().For(proj).Environments().For(env).Objects().Get(obj)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	if o == nil {
		rest.NotFoundResponse(w, r)
		return
	}
	rest.JSONResponse(w, r, o)
}
