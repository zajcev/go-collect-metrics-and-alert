package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	name  string
	value Metric
}

type Metric struct {
	mtype string
	value interface{}
}

var metrics = []MemStorage{}

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func metricCollector(w http.ResponseWriter, r *http.Request) {
	array := strings.Split(r.RequestURI, "/")
	if len(array) < 5 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		mt := array[2]
		mn := array[3]
		mv := array[4]
		if mt == "gauge" {
			v, err := parse(mv, mt)
			if err != nil {
				log.Default().Println(err)
				w.WriteHeader(http.StatusBadRequest)
			} else {
				metrics = append(metrics, MemStorage{mn, Metric{mt, v}})
				w.WriteHeader(http.StatusOK)
			}
		} else if mt == "counter" {
			v, err := parse(mv, mt)
			if err != nil {
				log.Default().Println(err)
				w.WriteHeader(http.StatusBadRequest)
			} else {
				metrics = append(metrics, MemStorage{mn, Metric{mt, v}})
				w.WriteHeader(http.StatusOK)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		fmt.Printf("%v \n", metrics)
	}
}

func parse(v string, t string) (interface{}, error) {
	if t == "gauge" {
		result, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else if t == "counter" {
		result, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		log.Printf("No metrics with type %v found", t)
	}
	return nil, nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, metricCollector)
	mux.HandleFunc(`/updater/`, metricCollector)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
