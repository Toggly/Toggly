package main

import (
	"context"
	"fmt"

	"github.com/Toggly/backend/app/rest/api"
)

//Application contains all internal components
type Application struct {
	api *api.TogglyApi
}

//Run application
func (a *Application) Run(ctx context.Context) error {

	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		fmt.Println("Done")
	}()

	return nil
}
