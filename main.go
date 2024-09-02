package main

import (
	"context"
	"net/http"
	"os"

	"github.com/DevJHansen/headlines/cmd"
	"github.com/DevJHansen/headlines/pkg/firebase"
)

func main() {
	// Get the PORT from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}

	ctx := context.Background()
	app, _ := firebase.NewFirebaseApp(ctx)

	http.HandleFunc("/scrape", func(w http.ResponseWriter, r *http.Request) {
		cmd.Handler(w, r, ctx, app)
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		// Write a "Hello, World!" response
		_, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			http.Error(w, "Unable to write response", http.StatusInternalServerError)
		}
	})

	// Start the server
	http.ListenAndServe(":"+port, nil)
}
