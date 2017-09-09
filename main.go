package main

import (
	"fmt"
	//	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	//	"github.com/tensorflow/tensorflow/tensorflow/go/op"
	//	"io/ioutil"
	//	"html/template"
	//	"log"
	//	"net/http"
	//"os"
)

/*
var indextmpl *template.Template
var resultstmpl *template.Template

func init() {
	index, err := Asset("pages/index.html")
	results, err := Asset("pages/results.html")

	if err != nil {
		log.Fatal(err)
	}

	indextmpl = template.Must(template.New("indextmpl").Parse(string(index)))
	resultstmpl = template.Must(template.New("resultstmpl").Parse(string(results)))

}*/

func main() {
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			indextmpl.Execute(w, map[string]string{"": ""})
			resultstmpl.Execute(w, map[string]string{"": ""})
		})

		log.Fatal(http.ListenAndServe(":8000", nil))
	*/
	fmt.Println(infer("./bird_mount_bluebird.jpg"))
	fmt.Println(labels[181])
}
