package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/api"
	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/go-chi/chi"
)

const (
	// ErrEnvironmentNotFound error
	ErrEnvironmentNotFound string = "Environment not found"
	// ErrEnvironmentNotEmpty error
	ErrEnvironmentNotEmpty string = "Environment not empty"
)

// EnvironmentCreateRequest type
type EnvironmentCreateRequest struct {
	Code        domain.EnvironmentCode
	Description string
	Protected   bool
}

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
		g.Post("/", a.createEnvironment)
		g.Put("/", a.updateEnvironment)
		g.Get("/{env_code}", a.cached(a.getEnvironment))
		g.Delete("/{env_code}", a.cached(a.deleteEnvironment))
	})
	return router
}

func (a *EnvironmentRestAPI) engine(r *http.Request) *api.EnvironmentAPI {
	return a.Engine.ForOwner(owner(r)).Projects().For(projectCode(r)).Environments()
}

func (a *EnvironmentRestAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return Cached(fn, a.Cache)
}

func (a *EnvironmentRestAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := a.engine(r).List()
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, "Project not found")
		default:
			log.Printf("[ERROR] %v", err)
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	response := map[string]interface{}{
		"environments": list,
	}
	JSONResponse(w, r, response)
}

func (a *EnvironmentRestAPI) getEnvironment(w http.ResponseWriter, r *http.Request) {
	env, err := a.engine(r).Get(environmentCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, ErrProjectNotFound)
		case api.ErrEnvironmentNotFound:
			NotFoundResponse(w, r, ErrEnvironmentNotFound)
		default:
			log.Printf("[ERROR] %v", err)
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, env)
}

func (a *EnvironmentRestAPI) deleteEnvironment(w http.ResponseWriter, r *http.Request) {
	err := a.engine(r).Delete(environmentCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, ErrProjectNotFound)
		case api.ErrEnvironmentNotFound:
			NotFoundResponse(w, r, ErrEnvironmentNotFound)
		case api.ErrEnvironmentNotEmpty:
			ErrorResponse(w, r, errors.New(ErrEnvironmentNotEmpty), http.StatusLocked)
		default:
			log.Printf("[ERROR] %v", err)
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, nil)
}

func (a *EnvironmentRestAPI) createEnvironment(w http.ResponseWriter, r *http.Request) {
	a.createUpdate(w, r, true)
}

func (a *EnvironmentRestAPI) updateEnvironment(w http.ResponseWriter, r *http.Request) {
	a.createUpdate(w, r, false)
}

func (a *EnvironmentRestAPI) createUpdate(w http.ResponseWriter, r *http.Request, create bool) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	env := &EnvironmentCreateRequest{}
	err = json.Unmarshal(body, env)
	if err != nil {
		ErrorResponse(w, r, errors.New("Bad request"), http.StatusBadRequest)
		return
	}
	var newEnv *domain.Environment
	if create {
		newEnv, err = a.engine(r).Create(env.Code, env.Description, env.Protected)
	} else {
		newEnv, err = a.engine(r).Update(env.Code, env.Description, env.Protected)
	}
	if err != nil {
		switch err {
		case api.ErrBadRequest:
			ErrorResponse(w, r, err, http.StatusBadRequest)
			return
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, ErrProjectNotFound)
			return
		case api.ErrEnvironmentNotFound:
			NotFoundResponse(w, r, ErrEnvironmentNotFound)
			return
		}

		switch err.(type) {
		case *storage.UniqueIndexError:
			ErrorResponse(w, r, err, http.StatusBadRequest)
		default:
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, newEnv)
}
