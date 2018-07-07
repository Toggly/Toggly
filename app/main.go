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

var revision = "development" //revision assigned on build

func printLogo() {
	fmt.Println("\n\x1b[94m" +
		"████████╗ ██████╗  ██████╗  ██████╗ ██╗  ██╗   ██╗\n" +
		"╚══██╔══╝██╔═══██╗██╔════╝ ██╔════╝ ██║  ╚██╗ ██╔╝\n" + "\x1b[34m" +
		"   ██║   ██║   ██║██║  ███╗██║  ███╗██║   ╚████╔╝ \n" +
		"   ██║   ██║   ██║██║   ██║██║   ██║██║    ╚██╔╝  \n" +
		"   ██║   ╚██████╔╝╚██████╔╝╚██████╔╝███████╗██║   \n" +
		"   ╚═╝    ╚═════╝  ╚═════╝  ╚═════╝ ╚══════╝╚═╝   \n" +
		"\x1b[0m\n")
	fmt.Printf("Toggly API Server (rev: %s)\n\n", revision)
	fmt.Println("--------------------------------------------------")
}

// Opts describes application command line arguments
type Opts struct {
	Port     int    `short:"p" long:"port" env:"TOGGLY_API_PORT" default:"8080" description:"port"`
	BasePath string `long:"base-path" env:"TOGGLY_API_BASE_PATH" default:"/api" description:"Base API Path"`
}

func main() {
	printLogo()
	var opts Opts
	if _, e := flags.NewParser(&opts, flags.Default).ParseArgs(os.Args[1:]); e != nil {
		os.Exit(1)
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
	app.Run(ctx)
	log.Println("[INFO] application terminated")
	log.Println("[INFO] Bye! 🖐")
}
