package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	flags "github.com/jessevdk/go-flags"
)

// Opts describes application command line arguments
type Opts struct {
	Port int `short:"p" long:"port" env:"TOGGLY_API_PORT" default:"8080" description:"port"`
}

var revision = "" //revision assigned on build

func main() {

	fmt.Printf("Toggly API Server %s\n", revision)

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	if _, e := p.ParseArgs(os.Args[1:]); e != nil {
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Print("[WARN] interrupt signal")
		cancel()
	}()

	app, err := createApplication(opts)
	if err != nil {
		log.Fatalf("[ERROR] failed to setup application, %+v", err)
	}

	err = app.Run(ctx)
	log.Printf("[INFO] application terminated %s", err)
}

func createApplication(opts Opts) (*Application, error) {
	return &Application{}, nil
}
