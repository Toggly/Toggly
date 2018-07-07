package main

import (
	"context"

	"github.com/Toggly/backend/app/cache"
	"github.com/Toggly/backend/app/rest/api"
)

//Application contains all internal components
type Application struct {
	api *api.TogglyAPI
}

//Run application
func (a *Application) Run(ctx context.Context) error {

	go func() {
		// stop on context cancellation
		<-ctx.Done()
		a.api.Stop()
	}()

	a.api.Run()

	return nil
}

func createApplication(opts Opts) (*Application, error) {
	var apiCache cache.DataCache
	var err error

	if apiCache, err = cache.NewHashMapCache(); err != nil {
		return nil, err
	}

	api := &api.TogglyAPI{
		Cache:    apiCache,
		BasePath: opts.BasePath,
		Port:     opts.Port,
	}

	return &Application{
		api: api,
	}, nil
}
