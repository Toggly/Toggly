package server

import (
	"context"

	"github.com/Toggly/core/internal/server/rest"
)

//Application contains all internal components
type Application struct {
	Router *rest.APIRouter
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
