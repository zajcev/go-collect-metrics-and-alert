package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"time"
)

var m Metrics
var counter int64 = 0

func monitor() {
	var rt runtime.MemStats
	runtime.ReadMemStats(&rt)
	mt := reflect.TypeOf(m)
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		setFieldValue(&m, f.Name, getValueByName(rt, f.Name))
		//log.Printf("Monitor: Name: %v = Value: %v", f.Name, getValueByName(m, f.Name))
	}
	addCustomMetric()
}

func reporter() {
	mt := reflect.TypeOf(m)
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		var t string
		baseUrl := "http://localhost:8080/update"
		v := getValueByName(m, f.Name)
		if reflect.TypeOf(v).String() == "float64" {
			t = "gauge"
		} else {
			t = "counter"
		}
		s := fmt.Sprintf("%v", v)
		res, err := url.JoinPath(baseUrl, t, f.Name, s)
		if err != nil {
			log.Printf("Error in reporter: %v", err)
		}
		resp, err := http.Post(res, "text/plain", nil)
		if err != nil {
			log.Printf("Error in reporter: %v", err)
		}
		resp.Close = true
		//log.Printf("Reporter: Name: %v = Value: %v", f.Name, getValueByName(&m, f.Name))
	}
}

func main() {
	monitorTimer := time.NewTicker(2 * time.Second)
	reporterTimer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-monitorTimer.C:
			monitor()
		case <-reporterTimer.C:
			reporter()
		}
	}
}
