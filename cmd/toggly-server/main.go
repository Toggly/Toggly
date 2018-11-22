package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Toggly/core/internal/server"
	"github.com/Toggly/core/internal/server/rest"

	"github.com/Toggly/core/internal/pkg/cache"
	"github.com/Toggly/core/internal/pkg/cache/cachedapi"
	"github.com/Toggly/core/internal/pkg/engine"
	"github.com/Toggly/core/internal/pkg/storage"
	"github.com/Toggly/core/internal/pkg/storage/mongo"
	flags "github.com/jessevdk/go-flags"
)

var revision = "development" //revision assigned on build

// Opts describes application command line arguments
type Opts struct {
	Toggly struct {
		Version  bool   `short:"v" long:"version"`
		Port     int    `short:"p" long:"port" env:"API_PORT" default:"8080" description:"Port"`
		BasePath string `long:"base-path" env:"API_BASE_PATH" default:"/api" description:"Base API Path"`
		NoLogo   bool   `long:"no-logo" description:"Do not show application logo"`
		Store    struct {
			Mongo struct {
				URL string `long:"url" env:"URL" description:"Mongo connection url"`
			} `group:"mongo" namespace:"mongo" env-namespace:"MONGO"`
		} `group:"store" namespace:"store" env-namespace:"STORE"`
		Cache struct {
			Type  string `long:"type" choice:"memory" choice:"redis" env:"TYPE" description:"Cache type"`
			Redis struct {
				URL string `long:"url" env:"URL" description:"Redis connection url"`
			} `group:"redis" namespace:"redis" env-namespace:"REDIS"`
		} `group:"cache" namespace:"cache" env-namespace:"CACHE"`
	} `group:"toggly" env-namespace:"TOGGLY"`
}

func main() {
	var dataStorage storage.DataStorage
	var dataCache cache.DataCache
	var err error

	var opts Opts
	if _, err = flags.NewParser(&opts, flags.Default).ParseArgs(os.Args[1:]); err != nil {
		os.Exit(0)
	}

	if opts.Toggly.Version {
		fmt.Printf("Toggly Server ver: %s\n", revision)
		os.Exit(0)
	}

	if !opts.Toggly.NoLogo {
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
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Print("[WARN] interrupt signal")
		cancel()
	}()

	switch opts.Toggly.Cache.Type {
	case "memory":
		dataCache = cache.NewInMemoryCache()
	case "redis":
		dataCache = cache.NewRedisCache(opts.Toggly.Cache.Redis.URL)
	default:
		log.Print("[WARN] No cache type specified. Cache disabled.")
	}

	if dataStorage, err = mongo.NewMongoStorage(opts.Toggly.Store.Mongo.URL); err != nil {
		log.Fatalf("[FATAL] Can't connect to storage: %+v", err)
	}

	server := &server.Application{
		Router: &rest.APIRouter{
			Version:  revision,
			API:      cachedapi.NewCachedAPI(engine.NewTogglyAPI(&dataStorage), dataCache),
			BasePath: opts.Toggly.BasePath,
			Port:     opts.Toggly.Port,
			IsDebug:  false,
		},
	}

	if err != nil {
		log.Fatalf("[FATAL] failed to setup application, %+v", err)
	}

	log.Print("[INFO] API server started")

	server.Run(ctx)
	log.Print("[INFO] application terminated")
	log.Print("[INFO] Bye!")
}

func centeredText(txt string, width int) string {
	return strings.Repeat(" ", (width-len(txt))/2) + txt
}
