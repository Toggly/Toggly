package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"plugin"
	"strings"
	"syscall"

	"github.com/Toggly/core/internal/server"
	"github.com/Toggly/core/internal/server/rest"

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
			Plugin struct {
				Name       string            `long:"name" env:"NAME" description:"Cache plugin name.\n Skip '-cache.so' suffix.\nFor example: '--cache.plugin.name=in-memory' will lookup 'in-memory-cache.so' file.\n"`
				Parameters map[string]string `long:"parameter" env:"PARAMETER" env-delim:"," description:"Plugin parameter.\nFor example: '--cache.plugin.parameter=param:value'.\n"`
			} `group:"plugin" namespace:"plugin" env-namespace:"PLUGIN"`
		} `group:"cache" namespace:"cache" env-namespace:"CACHE"`
	} `group:"toggly" env-namespace:"TOGGLY"`
}

func main() {
	var dataStorage storage.DataStorage
	var dataCache cachedapi.DataCache
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

	if opts.Toggly.Cache.Plugin.Name == "" {
		log.Print("[WARN] No cache plugin specified. Cache disabled.")
	} else {
		srcFile := opts.Toggly.Cache.Plugin.Name + "-cache.so"
		plug, err := plugin.Open(srcFile)
		if err != nil {
			log.Fatalf("Can't load plugin source `%s`: %v", srcFile, err)
		}
		getCache, err := plug.Lookup("GetCache")
		if err != nil {
			log.Fatalf("Can't lookup GetCache function in plugin source: %v", err)
		}
		fn := getCache.(func(map[string]string) interface {
			Get(key string) ([]byte, error)
			Set(key string, data []byte) error
			Flush(scopes ...string) error
		})
		dataCache = fn(opts.Toggly.Cache.Plugin.Parameters)
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
