package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ajorgensen/gochat/db"
	"github.com/ajorgensen/gochat/static"
	"github.com/ajorgensen/gochat/stream"
	"github.com/ajorgensen/gochat/templates"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	dbc, err := db.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbc.Close()

	if _, err := dbc.CreateConversation("foobar"); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServerFS(static.AssetsFS)))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := dbc.SelectConversations()

		if err != nil {
			log.Printf("err: %v", err)
			http.Error(w, "Error fetching conversations", http.StatusInternalServerError)
			return
		}

		if err := templates.Index(w, templates.IndexParams{Conversations: c}); err != nil {
			log.Printf("err: %v", err)
		}
	})

	r.Get("/c/{cID}", func(w http.ResponseWriter, r *http.Request) {
		cID := chi.URLParam(r, "cID")

		// Check to see if there is a conversation for this id
		conversation, err := dbc.FindConversation(cID)
		if err != nil {
			http.Error(w, "error fetching conversation", http.StatusInternalServerError)
			return
		}

		if conversation == nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}

		messages, err := dbc.GetMessages(conversation.ConversationID)
		if err != nil {
			http.Error(w, "error fetching messages", http.StatusInternalServerError)
			return
		}

		if err := templates.Chat(w, templates.ChatParams{
			Conversation: conversation,
			Messages:     messages,
		}); err != nil {
			log.Printf("err: %v", err)
		}
	})

	r.Post("/messages", messagesHandler(dbc))

	log.Printf("Listening on port 3333")
	http.ListenAndServe(":3333", r)
}

func messagesHandler(dbc *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		converationID := r.FormValue("conversationId")

		fmt.Printf("conversation id: %s\n", converationID)

		robotMessage := fmt.Sprintf("Hello, I am just a robot. Here is what you just said to me: %s", userMessage)

		stream := stream.New(robotMessage)
		words := stream.StreamWords()

		// Write the user message to the database
		err := dbc.CreateMessage(converationID, db.User, userMessage)
		if err != nil {
			http.Error(w, "error saving message", http.StatusInternalServerError)
			return
		}

		// Write robot message to the database
		err = dbc.CreateMessage(converationID, db.Assistant, robotMessage)
		if err != nil {
			http.Error(w, "error saving message", http.StatusInternalServerError)
			return
		}

		for word := range words {
			payload := &db.Message{
				ConversationID: converationID,
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
}
