package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const (
	docRoot       = "./static"
	indexTemplate = "index.html"
)

// Index is the front page of the web application
func (srv Web) Index(w http.ResponseWriter, r *http.Request) {
	t, err := srv.templates(indexTemplate)
	if err != nil {
		log.Printf("Error loading the application template: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (srv Web) templates(name string) (*template.Template, error) {
	// Parse the templates
	p := filepath.Join(docRoot, name)
	t, err := template.ParseFiles(p)
	if err != nil {
		log.Printf("Error loading the application template: %v\n", err)
	}
	return t, err
}
