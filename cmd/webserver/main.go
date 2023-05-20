package main

import (
	"log"
	"net/http"

	"github.com/brandonlbarrow/gonk/v2/internal/webserver"
)

func main() {
	m := http.NewServeMux()
	webserver.RegisterRoutes(m)
	if err := http.ListenAndServe(":8080", m); err != nil {
		log.Fatalf("error running webserver: %v", err)
	}
}
