package templates

import (
	"html/template"
	"io"
)

var t *template.Template

func init() {
	t = template.Must(template.ParseGlob("server/templates/*.html"))
}

func Render(w io.Writer, name string, data interface{}) error {
	return t.ExecuteTemplate(w, name, data)
}
