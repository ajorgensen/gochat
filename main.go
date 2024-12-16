package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/ajorgensen/gochat/static"
	"github.com/ajorgensen/gochat/stream"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gofrs/uuid"
)

type Conversation struct {
	ConversationID string    `json:"conversation_id"`
	Messages       []Message `json:"messages"`
}

type Message struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
}

var conversations = make(map[string]*Conversation)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServerFS(static.AssetsFS)))

	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", Conversation{})
	})
	r.Get("/c/{cID}", func(w http.ResponseWriter, r *http.Request) {
		cID := chi.URLParam(r, "cID")

		conversation := conversations[cID]
		if conversation == nil {
			http.Error(w, "Conversation not found", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "index.html", conversation)
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
	converationID := r.FormValue("conversation_id")
	fmt.Println(userMessage)

	var conversation *Conversation
	if converationID != "" {
		conversation = conversations[converationID]
	} else {
		conversation = &Conversation{
			ConversationID: uuid.Must(uuid.NewV4()).String(),
			Messages:       []Message{},
		}

		conversations[conversation.ConversationID] = conversation
	}

	stream := stream.New("Hello, I am just a robot")
	words := stream.StreamWords()

	// Append the users message
	conversation.Messages = append(conversation.Messages, Message{
		Message: userMessage,
	})

	// Append the robots message to the conversation
	conversation.Messages = append(conversation.Messages, Message{
		Message: stream.Message,
	})

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
