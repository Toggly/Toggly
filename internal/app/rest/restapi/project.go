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

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
)

// ProjectAPI servers project api namespace
type ProjectAPI struct {
	Cache  cache.DataCache
	Engine *api.Engine
}

// Routes returns routes for project namespace
func (p *ProjectAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Get("/", p.cached(p.list))
		group.Get("/{id}", p.cached(p.getProject))
		group.Post("/", p.saveProject)
	})
	return router
}

func (p *ProjectAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, p.Cache)
}

func (p *ProjectAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := p.Engine.ForOwner(rest.OwnerFromContext(r)).Project().List()
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
	id := chi.URLParam(r, "id")
	proj, err := p.Engine.ForOwner(rest.OwnerFromContext(r)).Project().Get(domain.ProjectCode(id))
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

func (p *ProjectAPI) saveProject(w http.ResponseWriter, r *http.Request) {
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
	err = p.Engine.ForOwner(rest.OwnerFromContext(r)).Project().Save(proj)
	if err != nil {
		switch err.(type) {
		case *storage.UniqueIndexError:
			rest.ErrorResponse(w, r, err, 400)
		default:
			rest.ErrorResponse(w, r, err, 500)
		}

	}
}
