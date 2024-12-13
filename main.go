package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/ajorgensen/gochat/static"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServerFS(static.AssetsFS)))

	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	r.Post("/messages", messagesHandler)

	log.Printf("Listening on port 3333")
	http.ListenAndServe(":3333", r)
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	r.ParseForm()
	userMessage := r.FormValue("message")

	// Immediately send a welcome message
	fmt.Fprintf(w, "data: Received your message: %s\n\n", userMessage)
	flusher.Flush()

	// Continuously send data every second until client disconnects
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-ctx.Done():
			// Client disconnected
			return
		case t := <-ticker.C:
			count++
			fmt.Fprintf(w, "data: %s - Count: %d\n\n", t.Format(time.RFC3339), count)
			flusher.Flush()
		}
	}
}
