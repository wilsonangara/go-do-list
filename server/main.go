package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kenshaw/envcfg"
)

const server = "backend"

var config *envcfg.Envcfg

func init() {
	var err error
	config, err = envcfg.New()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Printf("starting %s server", server)

	addr := net.JoinHostPort("", config.GetString("server.port"))
	srv := &http.Server{
		Addr:         addr,
		Handler:      registerHandlers(),
		ReadTimeout:  config.GetDuration("server.readTimeout"),
		IdleTimeout:  config.GetDuration("server.idleTimeout"),
		WriteTimeout: config.GetDuration("server.writeTimeout"),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown: %v", err)
		}
		log.Print("server shutdown")
		close(idleConnsClosed)
	}()

	log.Printf("starting %s on port: %s", server, addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server ListenAndServe: %v", err)
	}
	<-idleConnsClosed

	log.Print("server exited gracefully")
}

func registerHandlers() http.Handler {
	router := mux.NewRouter()

	return router
}
