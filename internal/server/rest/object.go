package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/go-chi/chi"
)

const (
	// ErrObjectNotFound error
	ErrObjectNotFound string = "Object not found"
	// ErrObjectHasInheritors error
	ErrObjectHasInheritors string = "Object has inheritors"
)

// ObjectCreateRequest type
type ObjectCreateRequest struct {
	Code        domain.ObjectCode
	Description string
	Inherits    *domain.ObjectInheritance
	Parameters  []*domain.Parameter
}

// ObjectRestAPI servers objects
type ObjectRestAPI struct {
	API api.TogglyAPI
}

// Routes returns routes for environments
func (a *ObjectRestAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Get("/", a.list)
		g.Post("/", a.createObject)
		g.Put("/", a.updateObject)
		g.Get("/{object_code}", a.getObject)
		g.Get("/{object_code}/inheritors", a.getObjectInheritors)
		g.Delete("/{object_code}", a.deleteObject)
	})
	return router
}

func (a *ObjectRestAPI) engine(r *http.Request) api.ObjectAPI {
	return a.API.ForOwner(owner(r)).Projects().For(projectCode(r)).Environments().For(environmentCode(r)).Objects()
}

func (a *ObjectRestAPI) list(w http.ResponseWriter, r *http.Request) {
	list, err := a.engine(r).List()
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, "Project not found")
		case api.ErrEnvironmentNotFound:
			NotFoundResponse(w, r, "Environment not found")
		default:
			log.Printf("[ERROR] %v", err)
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, list)
}

func (a *ObjectRestAPI) getObject(w http.ResponseWriter, r *http.Request) {
	obj, err := a.engine(r).Get(objectCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, ErrProjectNotFound)
		case api.ErrEnvironmentNotFound:
			NotFoundResponse(w, r, ErrEnvironmentNotFound)
		case api.ErrObjectNotFound:
			NotFoundResponse(w, r, ErrObjectNotFound)
		default:
			log.Printf("[ERROR] %v", err)
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, obj)
}

func (a *ObjectRestAPI) getObjectInheritors(w http.ResponseWriter, r *http.Request) {
	list, err := a.engine(r).InheritorsFlatList(objectCode(r))
	if err != nil {
		log.Printf("[ERROR] %v", err)
		ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, r, list)
}

func (a *ObjectRestAPI) deleteObject(w http.ResponseWriter, r *http.Request) {
	err := a.engine(r).Delete(objectCode(r))
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, ErrProjectNotFound)
		case api.ErrEnvironmentNotFound:
			NotFoundResponse(w, r, ErrEnvironmentNotFound)
		case api.ErrObjectNotFound:
			NotFoundResponse(w, r, ErrObjectNotFound)
		case api.ErrObjectHasInheritors:
			ErrorResponse(w, r, errors.New(ErrObjectHasInheritors), http.StatusLocked)
		default:
			log.Printf("[ERROR] %v", err)
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, nil)
}

func (a *ObjectRestAPI) createObject(w http.ResponseWriter, r *http.Request) {
	a.createUpdate(w, r, true)
}

func (a *ObjectRestAPI) updateObject(w http.ResponseWriter, r *http.Request) {
	a.createUpdate(w, r, false)
}

func (a *ObjectRestAPI) createUpdate(w http.ResponseWriter, r *http.Request, create bool) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	obj := &ObjectCreateRequest{}
	err = json.Unmarshal(body, obj)
	if err != nil {
		ErrorResponse(w, r, errors.New("Bad request"), http.StatusBadRequest)
		return
	}
	var newObj *domain.Object
	if create {
		newObj, err = a.engine(r).Create(&api.ObjectInfo{
			Code:        obj.Code,
			Description: obj.Description,
			Inherits:    obj.Inherits,
			Parameters:  obj.Parameters,
		})
	} else {
		newObj, err = a.engine(r).Update(&api.ObjectInfo{
			Code:        obj.Code,
			Description: obj.Description,
			Inherits:    obj.Inherits,
			Parameters:  obj.Parameters,
		})
	}
	if err != nil {
		switch err {
		case api.ErrProjectNotFound:
			NotFoundResponse(w, r, ErrProjectNotFound)
			return
		case api.ErrEnvironmentNotFound:
			NotFoundResponse(w, r, ErrEnvironmentNotFound)
			return
		case api.ErrObjectNotFound:
			NotFoundResponse(w, r, ErrObjectNotFound)
			return
		case api.ErrObjectParentNotExists:
			ErrorResponse(w, r, err, http.StatusBadRequest)
			return
		case api.ErrObjectInheritorTypeMismatch:
			ErrorResponse(w, r, err, http.StatusBadRequest)
			return
		}
		switch err.(type) {
		case *api.ErrBadRequest:
			ErrorResponse(w, r, err, http.StatusBadRequest)
			return
		case *storage.UniqueIndexError:
			ErrorResponse(w, r, err, http.StatusBadRequest)
		case *api.ErrObjectParameter:
			ErrorResponse(w, r, err, http.StatusBadRequest)
		default:
			ErrorResponse(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	JSONResponse(w, r, newObj)
}
