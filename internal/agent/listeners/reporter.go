package listeners

import (
	"bytes"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/cast"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"log"
	"net/http"
	"net/url"
	"reflect"
)

var MemStorage model.Metrics

func NewReporter(u string) {
	mt := reflect.TypeOf(MemStorage)
	for i := 0; i < mt.NumField(); i++ {
		mj := model.MetricJSON{}
		f := mt.Field(i)
		mj.ID = f.Name
		var t string
		u, err := url.Parse(u)
		if err != nil {
			log.Fatal(err)
		}
		v := model.GetValueByName(MemStorage, f.Name)
		if reflect.TypeOf(v).String() == "float64" {
			t = constants.Gauge
			result := cast.GetFloat(v)
			mj.Value = &result
		} else if reflect.TypeOf(v).String() == "int64" {
			t = constants.Counter
			result := cast.GetInt64(v)
			mj.Delta = &result
		}
		mj.MType = t
		req, err := json.Marshal(mj)
		if err != nil {
			log.Fatalf("Error marshalling json: %v", err)
		}

		resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(req))
		if err != nil {
			log.Printf("Error making POST request: %v", err)
			return
		}
		defer resp.Body.Close()
	}
}
