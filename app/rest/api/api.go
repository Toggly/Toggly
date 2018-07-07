package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Toggly/backend/app/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// TogglyAPI implements rest API
type TogglyAPI struct {
	Cache      cache.DataCache
	Port       int
	BasePath   string
	httpServer *http.Server
	lock       sync.Mutex
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
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Throttle(1000),
		middleware.Timeout(60*time.Second),
	)
	router.Route(a.BasePath, func(r chi.Router) {
		r.Mount("/v1", a.v1())
	})
	return router
}

func (a *TogglyAPI) v1() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Get("/project", a.getProject)
	})
	return router
}
