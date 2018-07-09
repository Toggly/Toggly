package val

import (
	"github.com/Toggly/backend/app/cache"
	"github.com/go-chi/chi"
)

// API servers value api namespace
type API struct {
	Cache cache.DataCache
}

// Routes returns routes for project namespace
func (a *API) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(g chi.Router) {
		g.Mount("/project", (&ProjectAPI{a.Cache}).Routes())
	})
	return router
}
