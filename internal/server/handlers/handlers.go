package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var metrics = models.NewMetricsStorage()

var htmlTemplate = `{{ range $key, $value := .Metrics}}
   <tr>Name: {{ $key }} Type: {{ .Type }} Value: {{ .Value }}</tr><br/>
{{ end }}`

// /update/{type}/{name}/{value}
func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	mname := chi.URLParam(r, "name")
	mtype := chi.URLParam(r, "type")
	mvalue := chi.URLParam(r, "value")
	if mtype == constants.Gauge {
		v, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			metrics.SetGauge(mname, mtype, v)
		}
	} else if mtype == constants.Counter {
		v, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			metrics.SetCounter(mname, mtype, v)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	err := r.Body.Close()
	if err != nil {
		log.Fatalf("Error while close body: %v", err)
	}
}

func GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	mname := chi.URLParam(r, "name")
	mtype := chi.URLParam(r, "type")
	g := metrics.GetMetric(mname, mtype)
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

func GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v", metrics)
	t := template.New("t")
	t, err := t.Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, metrics)
	if err != nil {
		panic("Failed to execute template: " + err.Error())
	}
	err = r.Body.Close()
	if err != nil {
		return
	}
}
