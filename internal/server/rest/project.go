package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/pkg/storage"

	"github.com/Toggly/core/internal/domain"

	"github.com/Toggly/core/internal/pkg/cache"

	"github.com/go-chi/chi"
)

const (
	// ErrProjectNotFound error
	ErrProjectNotFound string = "Project not found"
	// ErrProjectNotEmpty error
	ErrProjectNotEmpty string = "Project not empty"
)

// ProjectCreateRequest type
type ProjectCreateRequest struct {
	Code        domain.ProjectCode
	Description string
	Status      domain.ProjectStatus
}

// ProjectRestAPI servers project api namespace
type ProjectRestAPI struct {
	Cache cache.DataCache
	API   api.TogglyAPI
}

// Routes returns routes for project namespace
func (a *ProjectRestAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Get("/", a.list)
		group.Post("/", a.createProject)
		group.Put("/", a.updateProject)
		group.Get("/{project_code}", a.getProject)
		group.Delete("/{project_code}", a.deleteProject)
	})
	return router
}

func (a *ProjectRestAPI) engine(r *http.Request) api.ProjectAPI {
	return a.API.ForOwner(owner(r)).Projects()
}

func (a *ProjectRestAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := a.engine(r).List()
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
	proj, err := a.engine(r).Get(projectCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, ErrProjectNotFound)
		default:
			log.Printf("[ERROR] %v", err)
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, proj)
}

func (a *ProjectRestAPI) deleteProject(w http.ResponseWriter, r *http.Request) {
	err := a.engine(r).Delete(projectCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, ErrProjectNotFound)
		case api.ErrProjectNotEmpty:
			ErrorResponse(w, r, errors.New(ErrProjectNotEmpty), http.StatusLocked)
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
		p, err = a.engine(r).Create(&api.ProjectInfo{
			Code:        proj.Code,
			Description: proj.Description,
			Status:      proj.Status,
		})
	} else {
		p, err = a.engine(r).Update(&api.ProjectInfo{
			Code:        proj.Code,
			Description: proj.Description,
			Status:      proj.Status,
		})
	}
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, ErrProjectNotFound)
			return
		}
		switch err.(type) {
		case *api.ErrBadRequest:
			ErrorResponse(w, r, err, http.StatusBadRequest)
			return
		case *storage.UniqueIndexError:
			ErrorResponse(w, r, err, http.StatusBadRequest)
		default:
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, p)
}
