package main

import (
	//"fmt"
	//	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	//	"github.com/tensorflow/tensorflow/tensorflow/go/op"
	//	"io/ioutil"
	"html/template"
	"log"
	"net/http"
	//"os"
)

var tmpl *template.Template

func init() {
	data, err := Asset("pages/index.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl = template.Must(template.New("tmpl").Parse(string(data)))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		tmpl.Execute(w, map[string]string{"Name": "James"})
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
