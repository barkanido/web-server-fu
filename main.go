package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func run() error {
	var wait time.Duration
	var idealwait time.Duration
	defualtWait := time.Duration(15)
	defualtIdealWait := time.Duration(60)
	flag.DurationVar(&wait,
		"graceful-timeout",
		time.Second*defualtWait,
		"the duration for which the server gracefully wait"+
			"for existing connections to finish - e.g. 15s or 1m")
	flag.DurationVar(&idealwait,
		"ideal connection timeout",
		time.Second*defualtIdealWait,
		"ideal timeout")
	flag.Parse()
	// setup stuff
	database, err := setupDB()
	if err != nil {
		return fmt.Errorf("setup db: %w", err)
	}
	srv := newServer(database)
	httpSrv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: wait,
		ReadTimeout:  wait,
		IdleTimeout:  idealwait,
		Handler:      srv, // Pass our instance of gorilla/mux in.
	}
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	// Accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = httpSrv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down. err:", err)
	return nil
}

func setupDB() (*DB, error) {
	db := make(map[string]string)
	return &db, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
