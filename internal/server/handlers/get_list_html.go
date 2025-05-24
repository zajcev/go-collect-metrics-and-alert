package handlers

import (
	"context"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/config"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/db"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"html/template"
	"net/http"
	"time"
)

// var metrics = models.NewMetricsStorage()
var htmlTemplate = `{{ range $key, $value := .Metrics}}
   <tr>Name: {{ $key }} Type: {{ .MType }} Value: {{if .Delta}}{{.Delta}}{{end}} {{if .Value}}{{.Value}}{{end}}</tr><br/>
{{ end }}`

// GetAllMetrics return metrics in HTML
func GetAllMetrics(metrics *models.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.GetDBHost() != "" {
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()
			db.GetAllMetrics(ctx, metrics)
		}
		t := template.New("t")
		t, err := t.Parse(htmlTemplate)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "text/html")
		err = t.Execute(w, metrics)
		if err != nil {
			panic("Failed to execute template: " + err.Error())
		}
		err = r.Body.Close()
		if err != nil {
			return
		}
	}
}
