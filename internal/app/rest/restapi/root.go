package restapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Toggly/core/internal/app/api"
	"github.com/Toggly/core/internal/app/rest"
	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// APIRouter implements rest APIRouter
type APIRouter struct {
	Version    string
	Cache      cache.DataCache
	Engine     *api.Engine
	Port       int
	BasePath   string
	httpServer *http.Server
	lock       sync.Mutex
	IsDebug    bool
}

// Run rest api
func (api *APIRouter) Run() {
	routes := api.routes()
	api.lock.Lock()
	api.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", api.Port),
		Handler: chi.ServerBaseContext(context.Background(), routes),
	}
	api.lock.Unlock()
	log.Printf("[INFO] HTTP server listening on -> %s", api.httpServer.Addr)
	log.Printf("[INFO] APIRouter V.1 base path -> %s/v1", api.BasePath)
	err := api.httpServer.ListenAndServe()
	log.Printf("[INFO] HTTP server terminated, %s", err)
}

// Stop rest api
func (api *APIRouter) Stop() {
	log.Print("[INFO] stop REST server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	api.lock.Lock()
	if api.httpServer != nil {
		if err := api.httpServer.Shutdown(ctx); err != nil {
			log.Printf("[ERROR] REST stop error, %s", err)
		}
	}
	log.Print("[INFO] REST server stopped")
	api.lock.Unlock()
}

func (api *APIRouter) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Throttle(1000))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(rest.ServiceInfo("Toggly", api.Version))
	router.Route(api.BasePath, api.versions)
	if api.IsDebug {
		log.Println("[DEBUG] Profiler enabled on /debug path")
		router.Mount("/debug", middleware.Profiler())
	}
	return router
}

func (api *APIRouter) versions(r chi.Router) {
	r.Route("/v1", api.v1)
}

func (api *APIRouter) v1(r chi.Router) {
	r.Use(rest.AuthCtx)
	r.Use(rest.OwnerCtx)
	r.Use(rest.RequestIDCtx)
	r.Use(middleware.Logger)
	r.Use(rest.VersionCtx("v1"))
	r.Mount("/project", (&ProjectAPI{Cache: api.Cache, Engine: api.Engine}).Routes())
	// r.Mount("/project/{project_code}/env", (&EnvironmentAPIRouter{Cache: api.Cache, Storage: api.Storage}).Routes())
	// r.Mount("/project/{project_code}/env/{env_code}/object", (&ObjectAPIRouter{Cache: api.Cache, Storage: api.Storage}).Routes())
}
