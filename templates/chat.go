package templates

import (
	"io"

	"github.com/ajorgensen/gochat/db"
)

type ChatParams struct {
	Conversation *db.Conversation
	Messages     []*db.Message
	Providers    []Provider
}

func Chat(w io.Writer, p ChatParams) error {
	return chat.Execute(w, p)
}
