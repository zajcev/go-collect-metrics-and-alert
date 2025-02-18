package listeners

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		mj := model.MetricJson{}
		f := mt.Field(i)
		mj.ID = f.Name
		var t string
		u, err := url.Parse(u)
		if err != nil {
			log.Fatal(err)
		}
		p := "update"
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
		fmt.Printf("%+v\n", mj)
		res := u.JoinPath(p)
		req, err := json.Marshal(mj)
		_, err = http.Post(res.String(), "application/json", bytes.NewBuffer(req))
		//if err != nil {
		//	log.Printf("Error while request: %v", err)
		//}
		//log.Printf("Reporter: Name: %v = Value: %v", f.Name, model.GetValueByName(&MemStorage, f.Name))
	}
}
