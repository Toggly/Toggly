package restapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Toggly/core/internal/pkg/storage"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/pkg/api"
	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/server/rest"

	"github.com/go-chi/chi"
)

const (
	errProjectNotFound string = "Project not found"
)

// ProjectRestAPI servers project api namespace
type ProjectRestAPI struct {
	Cache  cache.DataCache
	Engine *api.Engine
}

// Routes returns routes for project namespace
func (a *ProjectRestAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Get("/", a.cached(a.list))
		group.Post("/", a.saveProject)
		group.Get("/{project_code}", a.cached(a.getProject))
		group.Delete("/{project_code}", a.cached(a.deleteProject))
	})
	return router
}

func (a *ProjectRestAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, a.Cache)
}

func (a *ProjectRestAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := a.Engine.ForOwner(owner(r)).Projects().List()
	if err != nil {
		log.Printf("[ERROR] %v", err)
		rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	rest.JSONResponse(w, r, list)
}

func (a *ProjectRestAPI) getProject(w http.ResponseWriter, r *http.Request) {
	proj, err := a.Engine.ForOwner(owner(r)).Projects().Get(projectCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			rest.NotFoundResponse(w, r, errProjectNotFound)
		default:
			log.Printf("[ERROR] %v", err)
			rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	if proj == nil {
		rest.NotFoundResponse(w, r, errProjectNotFound)
		return
	}
	rest.JSONResponse(w, r, proj)
}

func (a *ProjectRestAPI) deleteProject(w http.ResponseWriter, r *http.Request) {
	err := a.Engine.ForOwner(owner(r)).Projects().Delete(projectCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			rest.NotFoundResponse(w, r, errProjectNotFound)
		default:
			log.Printf("[ERROR] %v", err)
			rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	rest.JSONResponse(w, r, nil)
}

func (a *ProjectRestAPI) saveProject(w http.ResponseWriter, r *http.Request) {
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
	p, err := a.Engine.ForOwner(owner(r)).Projects().Create(proj.Code, proj.Description, proj.Status)
	if err != nil {
		switch err.(type) {
		case *storage.UniqueIndexError:
			rest.ErrorResponse(w, r, err, 400)
		default:
			rest.ErrorResponse(w, r, err, 500)
		}
		return
	}
	rest.JSONResponse(w, r, p)
}
