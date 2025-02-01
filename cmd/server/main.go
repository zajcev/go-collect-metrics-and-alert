package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

var metrics = map[string]*MemStorage{}

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func metricCollector(w http.ResponseWriter, r *http.Request) {
	array := strings.Split(r.RequestURI, "/")
	if len(array) < 5 {
		w.WriteHeader(http.StatusNotFound)
	} else if r.Method != "POST" || r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		mt, mn, mv := array[2], array[3], array[4]
		if mt == "gauge" {
			v, err := strconv.ParseFloat(mv, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(setGauge(metrics, mn, mt, v))
			}
		} else if mt == "counter" {
			v, err := strconv.ParseInt(mv, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(setCounter(metrics, mn, mt, v))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		err := r.Body.Close()
		if err != nil {
			log.Fatalf("Error while close body: %v", err)
		}
	}
	//for k, v := range metrics {
	//	fmt.Printf("key[%s] value[%s] \n -------------------------------------- \n", k, v)
	//}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, metricCollector)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
