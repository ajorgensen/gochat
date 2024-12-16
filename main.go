package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/ajorgensen/gochat/static"
	"github.com/ajorgensen/gochat/stream"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gofrs/uuid"
)

type Conversation struct {
	ConversationID string `json:"conversation_id"`
}

type Message struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
}

var conversation *Conversation

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

	//ctx := r.Context()
	r.ParseForm()
	userMessage := r.FormValue("message")
	fmt.Println(userMessage)

	// Continuously send data every second until client disconnects
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	if conversation == nil {
		conversation = &Conversation{
			ConversationID: uuid.Must(uuid.NewV4()).String(),
		}
	}

	stream := stream.New("Hello, I am just a robot")
	words := stream.StreamWords()

	for word := range words {
		payload := &Message{
			ConversationID: conversation.ConversationID,
			Message:        word,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, "data: %s\n\n", jsonPayload)
		flusher.Flush()
	}
}
