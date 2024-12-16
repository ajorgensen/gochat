package stream

import (
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

type Stream struct {
	Message string
}

func New(message string) *Stream {
	return &Stream{
		Message: message,
	}
}

func (s *Stream) StreamWords() chan string {
	wordChan := make(chan string)

	go func() {
		defer close(wordChan)

		// Split message into words
		words := strings.Fields(s.Message)
		var accumulator []string

		for _, word := range words {
			// Add new word to accumulator
			accumulator = append(accumulator, word)

			// Base delay of 300ms plus random jitter up to 700ms
			jitter := time.Duration(rand.Intn(300)) * time.Millisecond
			delay := 50*time.Millisecond + jitter

			time.Sleep(delay)
			// Send the accumulated string
			wordChan <- strings.Join(accumulator, " ")
		}
	}()

	return wordChan
}
