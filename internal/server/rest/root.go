package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/domain"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// APIRouter implements rest APIRouter
type APIRouter struct {
	Version    string
	API        api.TogglyAPI
	Port       int
	BasePath   string
	httpServer *http.Server
	lock       sync.Mutex
	IsDebug    bool
}

// Run rest api
func (r *APIRouter) Run() {
	routes := r.Router()
	r.lock.Lock()
	r.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", r.Port),
		Handler: chi.ServerBaseContext(context.Background(), routes),
	}
	r.lock.Unlock()
	log.Printf("[INFO] HTTP server listening on -> %s", r.httpServer.Addr)
	log.Printf("[INFO] APIRouter V.1 base path -> %s/v1", r.BasePath)
	err := r.httpServer.ListenAndServe()
	log.Printf("[INFO] HTTP server terminated, %s", err)
}

// Stop rest api
func (r *APIRouter) Stop() {
	log.Printf("[INFO] stop REST server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r.lock.Lock()
	if r.httpServer != nil {
		if err := r.httpServer.Shutdown(ctx); err != nil {
			log.Printf("[ERROR] REST stop error, %s", err)
		}
	}
	log.Print("[INFO] REST server stopped")
	r.lock.Unlock()
}

// Router returns router
func (r *APIRouter) Router() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Throttle(1000))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(ServiceInfo("Toggly", r.Version))
	router.Route(r.BasePath, r.versions)
	if r.IsDebug {
		log.Print("[DEBUG] Profiler enabled on /debug path")
		router.Mount("/debug", middleware.Profiler())
	}
	return router
}

func (r *APIRouter) versions(router chi.Router) {
	router.Route("/v1", r.v1)
}

func (r *APIRouter) v1(router chi.Router) {
	router.Use(middleware.Logger)
	router.Use(RequestIDCtx)
	router.Use(OwnerCtx)
	router.Use(VersionCtx("v1"))
	router.Mount("/project", (&ProjectRestAPI{API: r.API}).Routes())
	router.Mount("/project/{project_code}/env", (&EnvironmentRestAPI{API: r.API}).Routes())
	router.Mount("/project/{project_code}/env/{env_code}/object", (&ObjectRestAPI{API: r.API}).Routes())
}

func owner(r *http.Request) string {
	return OwnerFromContext(r)
}

func projectCode(r *http.Request) domain.ProjectCode {
	return domain.ProjectCode(chi.URLParam(r, "project_code"))
}

func environmentCode(r *http.Request) domain.EnvironmentCode {
	return domain.EnvironmentCode(chi.URLParam(r, "env_code"))
}

func objectCode(r *http.Request) domain.ObjectCode {
	return domain.ObjectCode(chi.URLParam(r, "object_code"))
}
