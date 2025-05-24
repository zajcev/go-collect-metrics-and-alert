package handlers

import (
	"encoding/json"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
)

// GetAllMetricsJSON return metrics in JSON
func GetAllMetricsJSON(metrics *models.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := metrics.GetAllMetrics()
		resp, err := json.Marshal(&m)
		if err != nil {
			log.Fatalf("Response: %v \n Error while writing response: %v", resp, err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		if err != nil {
			log.Fatalf("Error wgile write body : %v", err)
		}
	}
}
