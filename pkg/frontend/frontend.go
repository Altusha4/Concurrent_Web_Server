package frontend

import (
	_ "embed"
	"net/http"
	"strings"
)

//go:embed index.html
var indexHTML []byte

//go:embed style.css
var styleCSS []byte

//go:embed app.js
var appJS []byte

func ServeFrontend(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" || path == "/index.html" || (!strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/static/")) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
		return
	}

	if path == "/static/css/style.css" {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(styleCSS)
		return
	}

	if path == "/static/js/app.js" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(appJS)
		return
	}

	if strings.HasPrefix(path, "/static/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	ServeFrontend(w, r)
}
