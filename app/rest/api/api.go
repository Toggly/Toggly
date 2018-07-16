package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Toggly/core/app/storage"

	"github.com/Toggly/core/app/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// TogglyAPI implements rest API
type TogglyAPI struct {
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
func (a *TogglyAPI) Run() {
	router := a.routes()
	a.lock.Lock()
	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Port),
		Handler: router,
	}
	a.lock.Unlock()
	log.Printf("[INFO] HTTP server listening on → \x1b[1m%s\x1b[0m", a.httpServer.Addr)
	log.Printf("[INFO] API V.1 base path → \x1b[1m%s/v1\x1b[0m", a.BasePath)
	err := a.httpServer.ListenAndServe()
	log.Printf("[INFO] HTTP server terminated, %s", err)
}

// Stop rest api
func (a *TogglyAPI) Stop() {
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

func (a *TogglyAPI) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Throttle(1000))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(ServiceInfo("Toggly", a.Version))
	router.Route(a.BasePath, a.versions)
	if a.IsDebug {
		log.Println("[DEBUG] Profiler enabled on \x1b[1m/debug\x1b[0m path")
		router.Mount("/debug", middleware.Profiler())
	}
	return router
}

func (a *TogglyAPI) versions(r chi.Router) {
	r.Route("/v1", a.v1)
}

func (a *TogglyAPI) v1(r chi.Router) {
	r.Use(RequestIDCtx)
	r.Use(middleware.Logger)
	r.Use(VersionCtx("v1"))
	r.Mount("/project", (&ProjectAPI{Cache: a.Cache, Storage: a.Storage}).Routes())
	r.Mount("/project/{project_code}/object", (&ObjectAPI{Cache: a.Cache, Storage: a.Storage}).Routes())
	r.Mount("/project/{project_code}/env", (&EnvironmentAPI{Cache: a.Cache, Storage: a.Storage}).Routes())
	r.Mount("/project/{project_code}/env/{env_code}/object", (&ObjectAPI{Cache: a.Cache, Storage: a.Storage}).Routes())
}
