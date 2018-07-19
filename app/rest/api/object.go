package api

import (
	"log"
	"net/http"

	"github.com/Toggly/core/app/data"
	"github.com/Toggly/core/app/storage"

	"github.com/Toggly/core/app/rest"
	"github.com/Toggly/core/app/rest/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// ObjectAPI servers objects
type ObjectAPI struct {
	Cache   cache.DataCache
	Storage storage.DataStorage
}

// Routes returns routes for objects
func (p *ObjectAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Get("/", p.cached(p.list))
		g.Get("/{code}", p.cached(p.getObject))
	})
	return router
}

func (p *ObjectAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return cache.Cached(fn, p.Cache)
}

func (p *ObjectAPI) list(w http.ResponseWriter, r *http.Request) {
	proj := data.ProjectCode(chi.URLParam(r, "project_code"))
	env := data.EnvironmentCode(chi.URLParam(r, "env_code"))
	list, err := p.Storage.Projects(rest.CtxOwner(r)).For(proj).Environments().For(env).Objects().List()
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

func (p *ObjectAPI) getObject(w http.ResponseWriter, r *http.Request) {
	proj := data.ProjectCode(chi.URLParam(r, "project_code"))
	obj := data.ObjectCode(chi.URLParam(r, "code"))
	env := data.EnvironmentCode(chi.URLParam(r, "env_code"))
	o, err := p.Storage.Projects(rest.CtxOwner(r)).For(proj).Environments().For(env).Objects().Get(obj)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	if o == nil {
		rest.NotFoundResponse(w, r)
		return
	}
	render.JSON(w, r, o)
}
