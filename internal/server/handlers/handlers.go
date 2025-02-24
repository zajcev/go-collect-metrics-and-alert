package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var metrics = models.NewMetricsStorage()
var env = config.GetFlags()
var htmlTemplate = `{{ range $key, $value := .Metrics}}
   <tr>Name: {{ $key }} Type: {{ .MType }} Value: {{if .Delta}}{{.Delta}}{{end}} {{if .Value}}{{.Value}}{{end}}</tr><br/>
{{ end }}`

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
			syncWriter()
		}
	} else if mtype == constants.Counter {
		v, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			metrics.SetCounter(mname, mtype, v)
			syncWriter()
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	err := r.Body.Close()
	if err != nil {
		log.Fatalf("Error while close body: %v", err)
	}
}
func UpdateMetricHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		var m models.Metric
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			if m.MType == constants.Gauge {
				metrics.SetGaugeJSON(m)
			} else if m.MType == constants.Counter {
				metrics.SetCounterJSON(m)
			}
		}
		resp, err := json.Marshal(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		syncWriter()
	} else {
		w.WriteHeader(http.StatusBadRequest)
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

func GetMetricHandlerJSON(w http.ResponseWriter, r *http.Request) {
	var m models.Metric
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	value, code := metrics.GetMetricJSON(m)
	resp, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func GetAllMetrics(w http.ResponseWriter, r *http.Request) {
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

func GetAllMetricsJSON(w http.ResponseWriter, r *http.Request) {
	m := metrics.GetAllMetrics()
	resp, err := json.Marshal(&m)
	if err != nil {
		log.Fatalf("Response: %v \n Error while writing response: %v", resp, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func RestoreMetricStorage(file string) {
	consumer, err := storage.NewConsumer(file)
	if err != nil {
		log.Printf("Error while init file consumer %v", err)
		return
	}
	metrics, err = consumer.ReadMetrics()
	if err != nil {
		log.Printf("Error while read metric %v", err)
	}
}

func SaveMetricStorage(file string) {
	producer, err := storage.NewProducer(file)
	if err != nil {
		return
	}
	producer.WriteMetrics(metrics)
}

func syncWriter() {
	if env.StoreInterval == 0 {
		SaveMetricStorage(env.FilePath)
	}
}
