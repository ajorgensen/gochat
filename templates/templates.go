package templates

import (
	"embed"
	"html/template"
	"strings"
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

func parse(file string) *template.Template {
	return template.Must(
		template.New("layout.html").Funcs(funcs).ParseFS(tmplFS, "layout.html", file))
}
