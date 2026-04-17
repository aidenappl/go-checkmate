package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aidenappl/go-checkmate/env"
	"github.com/aidenappl/go-checkmate/middleware"
	"github.com/aidenappl/go-checkmate/poller"
	"github.com/aidenappl/go-checkmate/routers"
	"github.com/gorilla/mux"
)

func main() {

	p := poller.New(poller.LogSink{})
	defer p.Close()

	api := &routers.WorkerAPI{Poller: p}

	r := mux.NewRouter()

	// Healthcheck Endpoint
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Logging of requests
	r.Use(middleware.LoggingMiddleware)

	// 404 handler
	r.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Add Worker
	r.HandleFunc("/worker", api.NewWorkerHandler).Methods(http.MethodPost)

	// Users
	r.HandleFunc("/user", api.NewUserHandler).Methods(http.MethodPost)

	fmt.Println("✅ Checkmate API started on port", env.PORT)
	log.Fatal(http.ListenAndServe(":"+env.PORT, r))
}
