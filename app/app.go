package main

import (
	"context"

	"github.com/Toggly/core/app/cache"
	"github.com/Toggly/core/app/rest/api"
	"github.com/Toggly/core/app/storage"
)

//Application contains all internal components
type Application struct {
	api     *api.TogglyAPI
	storage *storage.DataStorage
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
	var dataStorage storage.DataStorage
	var err error
	if apiCache, err = cache.NewHashMapCache(); err != nil {
		return nil, err
	}
	if dataStorage, err = storage.NewFake(); err != nil {
		return nil, err
	}
	api := api.TogglyAPI{
		Cache:    apiCache,
		Storage:  dataStorage,
		BasePath: opts.BasePath,
		Port:     opts.Port,
	}
	return &Application{
		api: &api,
	}, nil
}
