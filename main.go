package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var indextmpl *template.Template
var resultstmpl *template.Template

var indexfile []byte
var resultsfile []byte

var err error

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type ResultsData struct {
	Token         string
	IsBird        bool
	ResultsString string
}

func init() {
	loadLabels()

	rand.Seed(time.Now().UnixNano())

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

func tokenGen() string { // Note: this function will fail after the year 2262, so it ain't my problem.
	b := make([]byte, 64)

	for i, cache, remain := 63, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		indextmpl.Execute(w, nil)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(handler.Header)
		defer file.Close()
	}
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {

	d := ResultsData{IsBird: false, ResultsString: ""}

	resultstmpl.Execute(w, d)
}

func main() {
	fmt.Println(tokenGen())

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/results", resultsHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))

	//	fmt.Println(infer("./bird_mount_bluebird.jpg"))
	//	fmt.Println(string(index))
}
