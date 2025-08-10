package middleware

import (
	"log"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("404 Not Found: %s %s", r.Method, r.RequestURI)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "Route not found"}`))
}
