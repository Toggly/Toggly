package server

import (
	"context"

	"github.com/Toggly/core/internal/server/rest/restapi"
)

//Application contains all internal components
type Application struct {
	Router *restapi.APIRouter
}

//Run application
func (a *Application) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		a.Router.Stop()
	}()
	a.Router.Run()
	return nil
}
