package templates

import (
	"io"

	"github.com/ajorgensen/gochat/db"
)

type IndexParams struct {
	Conversations []*db.Conversation
}

func Index(w io.Writer, p IndexParams) error {
	return index.Execute(w, p)
}
