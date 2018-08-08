package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Toggly/core/internal/app"
	"github.com/Toggly/core/internal/app/api"
	"github.com/Toggly/core/internal/app/rest/restapi"

	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/pkg/storage"
	flags "github.com/jessevdk/go-flags"
)

var revision = "development" //revision assigned on build

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
			Disabled bool `long:"disable" description:"Disable cache" env:"DISABLE"`
			Redis    struct {
				URL string `long:"url" env:"URL" description:"redis connection url"`
			} `group:"redis" namespace:"redis" env-namespace:"REDIS"`
		} `group:"cache" namespace:"cache" env-namespace:"CACHE"`
	} `group:"toggly" env-namespace:"TOGGLY"`
}

func main() {

	fmt.Println(`
::::::::::: ::::::::   ::::::::   ::::::::  :::     :::   ::: 
    :+:    :+:    :+: :+:    :+: :+:    :+: :+:     :+:   :+: 
    +:+    +:+    +:+ +:+        +:+        +:+      +:+ +:+  
    +#+    +#+    +:+ :#:        :#:        +#+       +#++:   
    +#+    +#+    +#+ +#+   +#+# +#+   +#+# +#+        +#+    
    #+#    #+#    #+# #+#    #+# #+#    #+# #+#        #+#    
    ###     ########   ########   ########  ########## ###    
	`)
	fmt.Println(centeredText("-= Core API Server =-", 63))
	fmt.Println(centeredText(fmt.Sprintf("ver: %s", revision), 63))
	fmt.Print("--------------------------------------------------------------\n\n")

	var apiCache cache.DataCache
	var dataStorage storage.DataStorage
	var err error

	var opts Opts
	if _, err = flags.NewParser(&opts, flags.Default).ParseArgs(os.Args[1:]); err != nil {
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

	if opts.Toggly.Cache.Disabled {
		log.Print("[WARN] \x1b[1mCACHE DISABLED\x1b[0m 😱")
	} else {
		if apiCache, err = cache.NewHashMapCache(); err != nil {
			log.Fatalf("Can't connect to cache service: %v", err)
		}
	}

	if dataStorage, err = storage.NewMongoStorage(opts.Toggly.Store.Mongo.URL); err != nil {
		log.Printf("Mongo URL: %s", opts.Toggly.Store.Mongo.URL)
		log.Fatalf("Can't connect to storeage: %v", err)
	}

	apiEngine := api.NewEngine(dataStorage)

	apiRouter := restapi.APIRouter{
		Version:  revision,
		Cache:    apiCache,
		Engine:   apiEngine,
		BasePath: opts.Toggly.BasePath,
		Port:     opts.Toggly.Port,
		IsDebug:  opts.Toggly.Debug,
	}

	app := &app.Application{
		Router: &apiRouter,
	}

	if err != nil {
		log.Fatalf("[ERROR] failed to setup application, %+v", err)
	}

	log.Print("[INFO] API server started \x1b[32m✔\x1b[0m")

	app.Run(ctx)
	log.Println("[INFO] application terminated")
	log.Println("[INFO] Bye! 🖐")
}

func centeredText(txt string, width int) string {
	return strings.Repeat(" ", (width-len(txt))/2) + txt
}
