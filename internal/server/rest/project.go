package rest

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

	"github.com/go-chi/chi"
)

const (
	errProjectNotFound string = "Project not found"
)

// ProjectCreateRequest type
type ProjectCreateRequest struct {
	Code        domain.ProjectCode
	Description string
	Status      domain.ProjectStatus
}

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
		group.Post("/", a.createProject)
		group.Put("/", a.updateProject)
		group.Get("/{project_code}", a.cached(a.getProject))
		group.Delete("/{project_code}", a.cached(a.deleteProject))
	})
	return router
}

func (a *ProjectRestAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return Cached(fn, a.Cache)
}

func (a *ProjectRestAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := a.Engine.ForOwner(owner(r)).Projects().List()
	if err != nil {
		log.Printf("[ERROR] %v", err)
		ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"projects": list,
	}
	JSONResponse(w, r, response)
}

func (a *ProjectRestAPI) getProject(w http.ResponseWriter, r *http.Request) {
	proj, err := a.Engine.ForOwner(owner(r)).Projects().Get(projectCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, errProjectNotFound)
		default:
			log.Printf("[ERROR] %v", err)
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, proj)
}

func (a *ProjectRestAPI) deleteProject(w http.ResponseWriter, r *http.Request) {
	err := a.Engine.ForOwner(owner(r)).Projects().Delete(projectCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, errProjectNotFound)
		default:
			log.Printf("[ERROR] %v", err)
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, nil)
}

func (a *ProjectRestAPI) createProject(w http.ResponseWriter, r *http.Request) {
	a.createUpdate(w, r, true)
}

func (a *ProjectRestAPI) updateProject(w http.ResponseWriter, r *http.Request) {
	a.createUpdate(w, r, false)
}

func (a *ProjectRestAPI) createUpdate(w http.ResponseWriter, r *http.Request, create bool) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	proj := &ProjectCreateRequest{}
	err = json.Unmarshal(body, proj)
	if err != nil {
		ErrorResponse(w, r, errors.New("Bad request"), http.StatusBadRequest)
		return
	}
	var p *domain.Project
	if create {
		p, err = a.Engine.ForOwner(owner(r)).Projects().Create(proj.Code, proj.Description, proj.Status)
	} else {
		p, err = a.Engine.ForOwner(owner(r)).Projects().Update(proj.Code, proj.Description, proj.Status)
		if err != nil && err == api.ErrProjectNotFound {
			NotFoundResponse(w, r, errProjectNotFound)
			return
		}
	}
	if err != nil {
		switch err.(type) {
		case *storage.UniqueIndexError:
			ErrorResponse(w, r, err, http.StatusBadRequest)
		default:
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, p)
}
