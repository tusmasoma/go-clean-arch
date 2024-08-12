package main

import (
	"net/http"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

func greet(_ http.ResponseWriter, _ *http.Request) {
	log.Info("Greet function is called")
}

func main() {
	http.HandleFunc("/", greet)
	if err := http.ListenAndServe(":8083", nil); err != nil { //nolint:gosec // This is just a sample code
		log.Info("Failed to start server", log.Ferror(err))
	}
}
