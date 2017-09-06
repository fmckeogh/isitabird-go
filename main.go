package main

import (
	//"fmt"
	//	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	//	"github.com/tensorflow/tensorflow/tensorflow/go/op"
	//	"io/ioutil"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

type Page struct {
	Title string
	Body  []byte
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("pages/" + tmpl + ".html")
	t.Execute(w, p)
}

func handler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "index", p)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
