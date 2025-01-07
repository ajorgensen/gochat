package templates

import (
	"embed"
	"html/template"
	"io"
	"strings"

	"github.com/ajorgensen/gochat/db"
)

//go:embed *.html
var tmplFS embed.FS

var funcs = template.FuncMap{
	"uppercase": func(v string) string {
		return strings.ToUpper(v)
	},
}

var (
	index = parse("index.html")
	chat  = parse("chat.html")
)

type IndexParams struct {
	Conversations []*db.Conversation
}

func Index(w io.Writer, p IndexParams) error {
	return index.Execute(w, p)
}

type ChatParams struct {
	Conversation *db.Conversation
	Messages     []*db.Message
}

func Chat(w io.Writer, p ChatParams) error {
	return chat.Execute(w, p)
}

func parse(file string) *template.Template {
	return template.Must(
		template.New("layout.html").Funcs(funcs).ParseFS(tmplFS, "layout.html", file))
}
