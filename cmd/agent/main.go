package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
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
	client := resty.New()
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		var t string
		u, err := url.Parse("http://localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
		p := "update"
		v := getValueByName(m, f.Name)
		if reflect.TypeOf(v).String() == "float64" {
			t = "gauge"
		} else {
			t = "counter"
		}
		s := fmt.Sprintf("%v", v)
		res := u.JoinPath(p, t, f.Name, s)
		resp, err := client.R().SetHeader("Content-Type", "text/plain").Post(res.String())
		if err != nil {
			log.Printf("Error while request: %v", err)
		}
		fmt.Println(resp)
		log.Printf("Reporter: Name: %v = Value: %v", f.Name, getValueByName(&m, f.Name))
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
