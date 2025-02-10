package listeners

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/agent/model"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"log"
	"net/url"
	"reflect"
)

var MemStorage model.Metrics
var counter int64 = 0

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
			t = constants.Gauge
		} else {
			t = constants.Counter
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
