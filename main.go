package main

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/poodbooq/pdf-meta-viewer/pdf"
)

var (
	port = ":8765"
	tmpl = new(template.Template)
)

const (
	tmplDir   = "static"
	tmplIndex = "index.html"
	tmplMeta  = "meta.html"
	tmplError = "error.html"
)

func init() {
	tmpl = template.Must(template.ParseFiles(
		filepath.Join(tmplDir, tmplIndex),
		filepath.Join(tmplDir, tmplMeta),
		filepath.Join(tmplDir, tmplError),
	))
}

func main() {
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		if err := tmpl.ExecuteTemplate(w, tmplIndex, port); err != nil {
			tmpl.ExecuteTemplate(w, tmplError, err)
		}

	case http.MethodPost:
		file, err := pdf.ReadPDFFromRequest(r)
		if err != nil {
			tmpl.ExecuteTemplate(w, tmplError, err)
			return
		}

		meta, err := file.GetMeta()
		if err != nil {
			tmpl.ExecuteTemplate(w, tmplError, err)
			return
		}

		if err := tmpl.ExecuteTemplate(w, tmplMeta, meta); err != nil {
			tmpl.ExecuteTemplate(w, tmplError, err)
			return
		}
	}
}
