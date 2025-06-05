package handlers

import (
	"context"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"html/template"
	"net/http"
	"time"
)

type GetMeticListStorage interface {
	GetAllMetrics(ctx context.Context) *models.MemStorage
}
type GetMetricListHandler struct {
	storage GetMeticListStorage
}

func NewGetMetricListHandler(storage GetMeticListStorage) *GetMetricListHandler {
	return &GetMetricListHandler{storage: storage}
}

var htmlTemplate = `{{ range $key, $value := .Storage}}
   <tr>Name: {{ $key }} Type: {{ .MType }} Value: {{if .Delta}}{{.Delta}}{{end}} {{if .Value}}{{.Value}}{{end}}</tr><br/>
{{ end }}`

// GetAllMetrics return metrics in HTML
func (handler *GetMetricListHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	metrics := handler.storage
	t := template.New("t")
	t, err := t.Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, metrics.GetAllMetrics(ctx))
	if err != nil {
		panic("Failed to execute template: " + err.Error())
	}
	err = r.Body.Close()
	if err != nil {
		return
	}
}
