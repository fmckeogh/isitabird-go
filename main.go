package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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

type Data struct {
	IsBird        bool
	ResultsString string
}

var Token string
var IsBird bool
var ResultsString string

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

func genToken() string { // Note: this function will fail after the year 2262, so it ain't my problem.
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
	//token := genToken()
	token := "testing"

	if r.Method == "GET" {
		indextmpl.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile(("uploadfile-" + token))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(handler.Header)

		defer file.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		probabilities, classes := infer(buf.Bytes())

		if classes[0] == 16 {
			IsBird = true
		} else {
			IsBird = false
		}

		var buffer bytes.Buffer

		for i := 0; i < 5; i++ {
			buffer.WriteString(" " + strconv.FormatFloat(float64(probabilities[i]*100), 'f', 2, 64) + "% ")
			buffer.WriteString((labels[int(classes[i])]))
			if i != 4 {
				buffer.WriteString(",")
			}
		}

		ResultsString = buffer.String()

		fmt.Println(ResultsString)

		http.Redirect(w, r, "/results", http.StatusSeeOther)
	}
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {

	d := Data{IsBird: IsBird, ResultsString: ResultsString}

	resultstmpl.Execute(w, d)
}

func main() {
	fmt.Println(genToken())

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/results", resultsHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
