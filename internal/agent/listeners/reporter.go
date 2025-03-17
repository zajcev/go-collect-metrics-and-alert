package listeners

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/convert"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

var MemStorage model.Metrics

func NewReporter(u string) {
	mt := reflect.TypeOf(MemStorage)
	mj := []model.MetricJSON{}
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		metric := model.MetricJSON{}
		metric.ID = f.Name
		var t string
		v := model.GetValueByName(MemStorage, f.Name)
		if reflect.TypeOf(v).String() == "float64" {
			t = constants.Gauge
			result := convert.GetFloat(v)
			metric.Value = &result
		} else if reflect.TypeOf(v).String() == "int64" {
			t = constants.Counter
			result := convert.GetInt64(v)
			metric.Delta = &result
		}
		metric.MType = t
		mj = append(mj, metric)
	}
	send(u, &mj)
}

func send(u string, list *[]model.MetricJSON) {
	req, err := json.Marshal(list)
	if err != nil {
		log.Fatalf("Error marshalling json: %v", err)
	}
	fu, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
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
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	request.Header.Add("Content-Encoding", "gzip")
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Type", "application/json")
	if config.GetHashKey() != "" {
		request.Header.Set("HashSHA256", calculateSHA256Hash(req, config.GetHashKey()))
	}
	resp, errs := retry(client, request, 3)
	if errs != nil {
		log.Printf("Error making request: %v", errs)
		return
	}
	defer resp.Body.Close()

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatalf("Error creating gzip reader: %v", err)
	}
	defer gzReader.Close()

	io.ReadAll(gzReader)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
}

func retry(client *http.Client, req *http.Request, count int) (*http.Response, error) {
	delay := 1
	var res error
	for i := 0; i < count; i++ {
		resp, err := client.Do(req)
		res = err
		if err == nil {
			return resp, nil
		}
		time.Sleep(time.Duration(delay) * time.Second)
		log.Printf("Retry after %v seconds", delay)
		delay += 2
	}
	return nil, res
}
func calculateSHA256Hash(data []byte, key string) string {
	k := []byte(key)
	signedData := append(k, data...)
	hash := sha256.Sum256(signedData)
	return hex.EncodeToString(hash[:])
}
