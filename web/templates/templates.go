package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html
var FS embed.FS

var T = template.Must(
	template.ParseFS(FS, "layout.html", "login.html"),
)
