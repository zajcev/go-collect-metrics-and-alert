package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var htmlTemplate = `{{ range $key, $value := . }}
   <tr>Name: {{ $key }} Type: {{ .Mtype }} Value: {{ .Value }}</tr><br/>
{{ end }}`

// /update/{type}/{name}/{value}
func updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	mname := chi.URLParam(r, "name")
	mtype := chi.URLParam(r, "type")
	mvalue := chi.URLParam(r, "value")
	if mtype == "gauge" {
		v, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			setGauge(metrics, mname, mtype, v)
		}
	} else if mtype == "counter" {
		v, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			setCounter(metrics, mname, mtype, v)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	err := r.Body.Close()
	if err != nil {
		log.Fatalf("Error while close body: %v", err)
	}
	for k, v := range metrics {
		fmt.Printf("key[%s] value[%s] \n -------------------------------------- \n", k, v)
	}
}

func getMetricHandler(w http.ResponseWriter, r *http.Request) {
	mname := chi.URLParam(r, "name")
	mtype := chi.URLParam(r, "type")
	g := getMetricValue(metrics, mname, mtype)
	if g != "" {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		res, err := w.Write([]byte(g))
		if err != nil {
			log.Fatalf("Response: %v \n Error while writing response: %v", res, err)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	err := r.Body.Close()
	if err != nil {
		return
	}
}

func getAllMetrics(w http.ResponseWriter, r *http.Request) {
	t := template.New("t")
	t, err := t.Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, metrics)
	if err != nil {
		panic(err)
	}
	err = r.Body.Close()
	if err != nil {
		return
	}
}
