package restapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Toggly/core/internal/app/rest"
	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// APIRouter implements rest APIRouter
type APIRouter struct {
	Version    string
	Cache      cache.DataCache
	Storage    storage.DataStorage
	Port       int
	BasePath   string
	httpServer *http.Server
	lock       sync.Mutex
	IsDebug    bool
}

// Run rest api
func (a *APIRouter) Run() {
	router := a.routes()
	a.lock.Lock()
	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Port),
		Handler: chi.ServerBaseContext(context.Background(), router),
	}
	a.lock.Unlock()
	log.Printf("[INFO] HTTP server listening on → \x1b[1m%s\x1b[0m", a.httpServer.Addr)
	log.Printf("[INFO] APIRouter V.1 base path → \x1b[1m%s/v1\x1b[0m", a.BasePath)
	err := a.httpServer.ListenAndServe()
	log.Printf("[INFO] HTTP server terminated, %s", err)
}

// Stop rest api
func (a *APIRouter) Stop() {
	log.Print("[INFO] stop REST server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a.lock.Lock()
	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			log.Printf("[ERROR] REST stop error, %s", err)
		}
	}
	log.Print("[INFO] REST server stopped")
	a.lock.Unlock()
}

func (a *APIRouter) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Throttle(1000))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(rest.ServiceInfo("Toggly", a.Version))
	router.Route(a.BasePath, a.versions)
	if a.IsDebug {
		log.Println("[DEBUG] Profiler enabled on \x1b[1m/debug\x1b[0m path")
		router.Mount("/debug", middleware.Profiler())
	}
	return router
}

func (a *APIRouter) versions(r chi.Router) {
	r.Route("/v1", a.v1)
}

func (a *APIRouter) v1(r chi.Router) {
	r.Use(rest.AuthCtx)
	r.Use(rest.OwnerCtx)
	r.Use(rest.RequestIDCtx)
	r.Use(middleware.Logger)
	r.Use(rest.VersionCtx("v1"))
	r.Mount("/project", (&ProjectAPI{Cache: a.Cache, Storage: a.Storage}).Routes())
	// r.Mount("/project/{project_code}/object", (&ObjectAPIRouter{Cache: a.Cache, Storage: a.Storage}).Routes())
	// r.Mount("/project/{project_code}/env", (&EnvironmentAPIRouter{Cache: a.Cache, Storage: a.Storage}).Routes())
	// r.Mount("/project/{project_code}/env/{env_code}/object", (&ObjectAPIRouter{Cache: a.Cache, Storage: a.Storage}).Routes())
}
