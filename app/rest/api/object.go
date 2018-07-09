package api

import (
	"log"
	"net/http"

	"github.com/Toggly/backend/app/data"

	"github.com/Toggly/backend/app/storage"

	"github.com/Toggly/backend/app/cache"
	"github.com/Toggly/backend/app/rest"
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

func (p *ObjectAPI) cached(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return cache.Cached(fn, p.Cache)
}

func (p *ObjectAPI) list(w http.ResponseWriter, r *http.Request) {
	proj := data.CodeType(chi.URLParam(r, "project_code"))
	list, err := p.Storage.ListObjects(proj)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, list)
}

func (p *ObjectAPI) getObject(w http.ResponseWriter, r *http.Request) {
	proj := data.CodeType(chi.URLParam(r, "project_code"))
	obj := data.CodeType(chi.URLParam(r, "code"))
	o, err := p.Storage.GetObject(proj, obj)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, o)
}
