package listeners

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
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
		fu, err := url.Parse(u)
		if err != nil {
			log.Fatal(err)
		}
		v := model.GetValueByName(MemStorage, f.Name)
		if reflect.TypeOf(v).String() == "float64" {
			t = constants.Gauge
			result := convert.GetFloat(v)
			mj.Value = &result
		} else if reflect.TypeOf(v).String() == "int64" {
			t = constants.Counter
			result := convert.GetInt64(v)
			mj.Delta = &result
		}
		mj.MType = t
		req, err := json.Marshal(mj)
		if err != nil {
			log.Fatalf("Error marshalling json: %v", err)
		}

		var buf bytes.Buffer
		g := gzip.NewWriter(&buf)
		if _, err = g.Write(req); err != nil {
			log.Fatalf("Error compressing json: %v", err)
			return
		}
		if err = g.Close(); err != nil {
			log.Fatalf("Error compressing json: %v", err)
			return
		}

		client := &http.Client{}
		request, err := http.NewRequest("POST", fu.String(), &buf)
		request.Header.Add("Content-Encoding", "gzip")
		request.Header.Add("Accept-Encoding", "gzip")
		request.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(request)
		if err != nil {
			log.Printf("Error making request: %v", err)
			return
		}
		defer resp.Body.Close()
	}
}
