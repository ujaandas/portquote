package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html
var FS embed.FS

var LoginT = template.Must(
	template.ParseFS(FS, "layout.html", "login.html"),
)

var DashT = template.Must(
	template.ParseFS(FS, "layout.html", "dashboard.html"),
)
