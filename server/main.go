package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/kenshaw/envcfg"
	"github.com/twitchtv/twirp"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/twitchtv/twirp"
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
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("server listening to %v", l.Addr())
	}()

	// catch interruption signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	// graceful shutdown
	shutDownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	log.Print("server exited gracefully")
}

func setupTwirpServer(*twirp.WrapServer) {

}
