package ui

import (
	"fmt"
	"log"
	"net/http"
)

// Start will start the HTTP server on the specified port, defaulting to 3000 if not provided
func Start(port string) {
	if port == "" {
		port = "8888" // Set default port to 3000
	}

	http.HandleFunc("/", homePageHandler)

	fmt.Printf("Server listening on port %s\n", port)
	log.Panic(http.ListenAndServe(":"+port, nil))
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "hello GovCMS!")
	if err != nil {
		log.Panic(err)
	}
}
