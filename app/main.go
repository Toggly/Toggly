package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Toggly/core/app/rest/api"
	"github.com/Toggly/core/app/rest/cache"
	"github.com/Toggly/core/app/storage"
	flags "github.com/jessevdk/go-flags"
)

var revision = "development" //revision assigned on build

func printLogo() {
	fmt.Println(`
::::::::::: ::::::::   ::::::::   ::::::::  :::     :::   ::: 
    :+:    :+:    :+: :+:    :+: :+:    :+: :+:     :+:   :+: 
    +:+    +:+    +:+ +:+        +:+        +:+      +:+ +:+  
    +#+    +#+    +:+ :#:        :#:        +#+       +#++:   
    +#+    +#+    +#+ +#+   +#+# +#+   +#+# +#+        +#+    
    #+#    #+#    #+# #+#    #+# #+#    #+# #+#        #+#    
    ###     ########   ########   ########  ########## ###    
	`)
	fmt.Println(centered("-= Core API Server =-", 63))
	fmt.Println(centered(fmt.Sprintf("ver: %s", revision), 63))
	fmt.Print("--------------------------------------------------------------\n\n")
}

func centered(txt string, width int) string {
	return strings.Repeat(" ", (width-len(txt))/2) + txt
}

// Opts describes application command line arguments
type Opts struct {
	Toggly struct {
		Port     int    `short:"p" long:"port" env:"API_PORT" default:"8080" description:"port"`
		BasePath string `long:"base-path" env:"API_BASE_PATH" default:"/api" description:"Base API Path"`
		Debug    bool   `long:"debug" description:"Run in DEBUG mode"`
		Store    struct {
			Mongo struct {
				URL string `long:"url" env:"URL" description:"mongo connection url"`
			} `group:"mongo" namespace:"mongo" env-namespace:"MONGO"`
		} `group:"store" namespace:"store" env-namespace:"STORE"`
		Cache struct {
			Disabled bool `long:"disable" description:"Disable cache"`
			Redis    struct {
				URL string `long:"url" env:"URL" description:"redis connection url"`
			} `group:"redis" namespace:"redis" env-namespace:"REDIS"`
		} `group:"cache" namespace:"cache" env-namespace:"CACHE"`
	} `group:"toggly" env-namespace:"TOGGLY"`
}

func main() {
	printLogo()
	var opts Opts
	if _, e := flags.NewParser(&opts, flags.Default).ParseArgs(os.Args[1:]); e != nil {
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Print("[WARN] interrupt signal \x1b[31m✘\x1b[0m")
		cancel()
	}()

	app, err := createApplication(opts)
	if err != nil {
		log.Fatalf("[ERROR] failed to setup application, %+v", err)
	}

	log.Print("[INFO] API server started \x1b[32m✔\x1b[0m")
	if opts.Toggly.Cache.Disabled {
		log.Print("[WARN] \x1b[1mCache disabled\x1b[0m")
	}
	app.Run(ctx)
	log.Println("[INFO] application terminated")
	log.Println("[INFO] Bye! 🖐")
}

//Application contains all internal components
type Application struct {
	api     *api.TogglyAPI
	storage *storage.DataStorage
}

//Run application
func (a *Application) Run(ctx context.Context) error {
	go func() {
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
	if apiCache, err = cache.NewHashMapCache(!opts.Toggly.Cache.Disabled); err != nil {
		return nil, err
	}
	if dataStorage, err = storage.NewMongoStorage(opts.Toggly.Store.Mongo.URL); err != nil {
		return nil, err
	}
	api := api.TogglyAPI{
		Version:  revision,
		Cache:    apiCache,
		Storage:  dataStorage,
		BasePath: opts.Toggly.BasePath,
		Port:     opts.Toggly.Port,
		IsDebug:  opts.Toggly.Debug,
	}
	return &Application{
		api: &api,
	}, nil
}
