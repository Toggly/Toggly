package restapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Toggly/core/internal/pkg/storage"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/app/api"
	"github.com/Toggly/core/internal/app/rest"
	"github.com/Toggly/core/internal/pkg/cache"

	"github.com/go-chi/chi"
)

// ProjectAPI servers project api namespace
type ProjectAPI struct {
	Cache  cache.DataCache
	Engine *api.Engine
}

// Routes returns routes for project namespace
func (api *ProjectAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Get("/", api.cached(api.list))
		group.Get("/{id}", api.cached(api.getProject))
		group.Post("/", api.saveProject)
	})
	return router
}

func (api *ProjectAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, api.Cache)
}

func (api *ProjectAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := api.Engine.ForOwner(rest.OwnerFromContext(r)).Project().List()
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

func (api *ProjectAPI) getProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	proj, err := api.Engine.ForOwner(rest.OwnerFromContext(r)).Project().Get(domain.ProjectCode(id))
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	if proj == nil {
		rest.NotFoundResponse(w, r)
		return
	}
	rest.JSONResponse(w, r, proj)
}

func (api *ProjectAPI) saveProject(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.ErrorResponse(w, r, err, 500)
		return
	}
	proj := &domain.Project{}
	err = json.Unmarshal(body, proj)
	if err != nil {
		rest.ErrorResponse(w, r, errors.New("Bad request"), 400)
		return
	}
	err = api.Engine.ForOwner(rest.OwnerFromContext(r)).Project().Save(proj)
	if err != nil {
		switch err.(type) {
		case *storage.UniqueIndexError:
			rest.ErrorResponse(w, r, err, 400)
		default:
			rest.ErrorResponse(w, r, err, 500)
		}
	}
}
