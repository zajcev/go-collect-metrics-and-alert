package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	db          *DatabaseStorage
	pgContainer testcontainers.Container
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		),
	}

	var err error
	pgContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("Failed to start container: %v\n", err)
		os.Exit(1)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		fmt.Printf("Failed to get container host: %v\n", err)
		os.Exit(1)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Printf("Failed to get container port: %v\n", err)
		os.Exit(1)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=testuser password=testpass dbname=testdb sslmode=disable",
		host, port.Int())
	pool, err := pgxpool.New(ctx, connStr)

	var maxAttempts = 5
	for i := 0; i < maxAttempts; i++ {
		db, err = NewDatabaseStorage(pool)
		if i < maxAttempts-1 {
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	Migration(connStr, "scripts")
	code := m.Run()

	if err = pgContainer.Terminate(ctx); err != nil {
		fmt.Printf("Failed to terminate container: %v\n", err)
	}
	os.Exit(code)
}

func TestPing(t *testing.T) {
	err := db.Ping(context.Background())
	if err != nil {
		t.Fatalf("Error ping database : %v", err)
	}
}

func TestSetDeltaRaw(t *testing.T) {
	db.SetDeltaRaw(context.Background(), "testDelta", "gauge", 30.0)
	testDeltaValue := db.GetMetricRaw(context.Background(), "testDelta", "gauge")

	assert.Equal(t, int64(30), testDeltaValue)
}

func TestSetValueRaw(t *testing.T) {
	db.SetValueRaw(context.Background(), "testValue", "counter", 30)
	testCounterValue := db.GetMetricRaw(context.Background(), "testValue", "counter")
	assert.Equal(t, float64(30), testCounterValue)
}

func TestSetDeltaJSON(t *testing.T) {
	delta := int64(30)
	metric := models.Metric{ID: "testCounter", MType: "counter", Delta: &delta}

	db.SetDeltaJSON(context.Background(), metric)
	result, _ := db.GetMetricJSON(context.Background(), metric)

	if *result.Delta != delta {
		t.Errorf("Expected Delta %d, got %d", delta, *result.Delta)
	}
}

func TestSetValueJSON(t *testing.T) {
	value := float64(30)
	metric := models.Metric{ID: "testGauge", MType: "gauge", Value: &value}

	db.SetValueJSON(context.Background(), metric)
	result, _ := db.GetMetricJSON(context.Background(), metric)

	if *result.Value != value {
		t.Errorf("Expected Delta %v, got %v", value, *result.Value)
	}
}

func TestSetListJSON(t *testing.T) {
	value := float64(123.45)
	delta := int64(42)
	list := []models.Metric{
		{
			ID:    "alloc",
			MType: "gauge",
			Value: &value,
		},
		{
			ID:    "pollCount",
			MType: "counter",
			Delta: &delta,
		},
	}
	code := db.SetListJSON(context.Background(), list)
	result := db.GetAllMetrics(context.Background())
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("JSON marshaling failed: %v", err)
	}
	t.Logf("Result JSON: %s", string(jsonData))
	assert.Equal(t, 200, code)
}
