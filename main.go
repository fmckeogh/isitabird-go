package main

import (
	"fmt"
	"html/template"
	"log"

	_ "errors"
	_ "io/ioutil"
	"net/http"
	_ "os"
)

var indextmpl *template.Template
var resultstmpl *template.Template

var indexfile []byte
var resultsfile []byte

var err error

type ResultsData struct {
	IsBird        bool
	ResultsString string
}

func init() {
	loadLabels()

	indexfile, err = Asset("pages/index.html")
	if err != nil {
		log.Fatal(err)
	}

	resultsfile, err = Asset("pages/results.html")
	if err != nil {
		log.Fatal(err)
	}

	indextmpl = template.Must(template.New("indextmpl").Parse(string(indexfile)))
	resultstmpl = template.Must(template.New("resultstmpl").Parse(string(resultsfile)))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	indextmpl.Execute(w, nil)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)

	d := ResultsData{IsBird: false, ResultsString: "523%%%%dgsgsdgw       &&&&,,,,,,&7,7,7jksndfhakfnadtever"}

	resultstmpl.Execute(w, d)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/results", resultsHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))

	//	fmt.Println(infer("./bird_mount_bluebird.jpg"))
	//	fmt.Println(string(index))
}
