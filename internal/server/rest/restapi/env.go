package restapi

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
	"github.com/Toggly/core/internal/server/rest"
	"github.com/go-chi/chi"
)

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
		g.Post("/", a.saveEnv)
		// g.Get("/{code}", a.cached(a.getEnvironment))
	})
	return router
}

func (a *EnvironmentRestAPI) cached(fn http.HandlerFunc) http.HandlerFunc {
	return rest.Cached(fn, a.Cache)
}

func (a *EnvironmentRestAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := a.Engine.ForOwner(owner(r)).Projects().For(projectCode(r)).Environments().List()
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			rest.NotFoundResponse(w, r, "Project not found")
		default:
			log.Printf("[ERROR] %v", err)
			rest.ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	rest.JSONResponse(w, r, list)
}

func (a *EnvironmentRestAPI) saveEnv(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.ErrorResponse(w, r, err, 500)
		return
	}
	env := &domain.Environment{}
	err = json.Unmarshal(body, env)
	if err != nil {
		rest.ErrorResponse(w, r, errors.New("Bad request"), 400)
		return
	}
	e, err := a.Engine.ForOwner(owner(r)).Projects().For(projectCode(r)).Environments().Create(env.Code, env.Description, env.Protected)
	if err != nil {
		switch err.(type) {
		case *storage.UniqueIndexError:
			rest.ErrorResponse(w, r, err, 400)
		default:
			rest.ErrorResponse(w, r, err, 500)
		}
		return
	}
	rest.JSONResponse(w, r, e)
}
