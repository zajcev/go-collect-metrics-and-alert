package listeners

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"log"
	"math/rand"
	"net/url"
	"reflect"
	"runtime"
)

var MemStorage model.Metrics
var counter int64 = 0

func NewMonitor() {
	var rt runtime.MemStats
	runtime.ReadMemStats(&rt)
	mt := reflect.TypeOf(MemStorage)
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		model.SetFieldValue(&MemStorage, f.Name, model.GetValueByName(rt, f.Name))
		//log.Printf("Monitor: Name: %v = Value: %v", f.Name, getValueByName(m, f.Name))
	}
	AddCustomMetric()
}

func NewReporter(u string) {
	mt := reflect.TypeOf(MemStorage)
	client := resty.New()
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		var t string
		u, err := url.Parse(u)
		if err != nil {
			log.Fatal(err)
		}
		p := "update"
		v := model.GetValueByName(MemStorage, f.Name)
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
		log.Printf("Reporter: Name: %v = Value: %v", f.Name, model.GetValueByName(&MemStorage, f.Name))
	}
}

func AddCustomMetric() {
	model.SetFieldValue(&MemStorage, "RandomValue", rand.Float64())
	if model.GetValueByName(MemStorage, "PollCount") == nil {
		model.SetFieldValue(MemStorage, "PollCount", int64(1))
	} else {
		counter = model.GetValueByName(&MemStorage, "PollCount").(int64)
		counter++
		model.SetFieldValue(&MemStorage, "PollCount", counter)
	}
}
