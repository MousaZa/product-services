package main

import (
	"context"
	"github.com/MousaZa/product-services/product-images/files"
	"github.com/MousaZa/product-services/product-images/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "/filestore/tmp/images", "Base path to save images")

func main() {

	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
			Level: hclog.LevelFromString(*logLevel),
		},
	)

	// create a logger for the server from the default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the storage class, use local storage
	// max filesize 5MB
	stor, err := files.NewLocal(*basePath, 1024*1000*5)
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	// create the handlers
	fh := handlers.NewFiles(stor, l)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	mw := handlers.GzipHandler{}
	// problem with FileServer is that it is dumb
	getRouter.Handle("/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", http.FileServer(http.Dir(*basePath)))
	getRouter.Use(mw.GzipMiddleware)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.UploadRest)
	postRouter.HandleFunc("/", fh.UploadMultipart)
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}),
		gohandlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete}))
	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	// start the server
	go func() {
		l.Info("Starting server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
