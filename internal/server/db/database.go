package db

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
)

type Database interface {
	SetDeltaRaw(ctx context.Context, name string, metricType string, value int64) int
	SetValueRaw(ctx context.Context, name string, metricType string, value float64) int
	SetValueJSON(ctx context.Context, input models.Metric) int
	SetDeltaJSON(ctx context.Context, input models.Metric) int
	GetMetricRaw(ctx context.Context, name string, metricType string) any
	GetMetricJSON(ctx context.Context, input models.Metric) (models.Metric, int)
	GetAllMetrics(ctx context.Context, ms *models.MemStorage) *models.MemStorage
	SetListJSON(ctx context.Context, list []models.Metric) int
	Ping(ctx context.Context) error
}
type DatabaseStorage struct {
	Database *pgxpool.Pool
}

func NewDatabaseStorage(database *pgxpool.Pool) (*DatabaseStorage, error) {
	return &DatabaseStorage{
		Database: database,
	}, nil
}

//go:embed scripts/000001_test.up.sql
var file string

// migration executes database migration
func Migration(DBUrl string, path string) {
	gp, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error get current directory : %v", err)
	}
	filePath := filepath.Join(gp, path)
	d, _ := sql.Open("postgres", DBUrl)
	driver, _ := postgres.WithInstance(d, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///"+filePath, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err = m.Up(); err != nil {
		log.Printf("Migrating database: %v", err)
	}
}

// Ping checks database connection
func (ds *DatabaseStorage) Ping(ctx context.Context) error {
	return ds.Database.Ping(ctx)
}

// GetMetricRaw returns metric value from database
func (ds *DatabaseStorage) GetMetricRaw(ctx context.Context, name string, metricType string) any {
	var value interface{}
	if metricType == constants.Gauge {
		err := ds.Database.QueryRow(ctx, getDelta, name, metricType).Scan(&value)
		if err != nil {
			return nil
		}
		return value
	} else {
		err := ds.Database.QueryRow(ctx, getValue, name, metricType).Scan(&value)
		if err != nil {
			return nil
		}
		return value
	}
}

// GetMetricJSON returns models.Metric for response in JSON format
func (ds *DatabaseStorage) GetMetricJSON(ctx context.Context, m models.Metric) (models.Metric, int) {
	row, _ := ds.Database.Query(ctx, getMetric, m.ID, m.MType)
	if row.Err() != nil {
		log.Printf("Error while execute query: %v", row.Err())
	}
	defer row.Close()
	if row.Next() {
		err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			return models.Metric{}, http.StatusInternalServerError
		}
		return m, http.StatusOK
	}
	return models.Metric{}, http.StatusNotFound
}

// SetDeltaRaw sets value for metric type Gauge in a database by raw value
func (ds *DatabaseStorage) SetDeltaRaw(ctx context.Context, name string, metricType string, value int64) int {
	_, err := ds.Database.Exec(ctx, insertDelta, name, metricType, value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, delta=%v", err, name, metricType, value)
		}
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// SetValueRaw sets value for metric type Counter in a database by raw value
func (ds *DatabaseStorage) SetValueRaw(ctx context.Context, name string, metricType string, value float64) int {
	_, err := ds.Database.Exec(ctx, insertValue, name, metricType, value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, value=%v", err, name, metricType, value)
		}
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// SetDeltaJSON sets value for metric type Gauge in a database by JSON value
func (ds *DatabaseStorage) SetDeltaJSON(ctx context.Context, m models.Metric) int {
	_, err := ds.Database.Exec(ctx, insertDelta, m.ID, m.MType, m.Delta)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, delta=%v", err, m.ID, m.MType, m.Delta)
		}
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// SetValueJSON sets value for metric type Counter in a database by JSON value
func (ds *DatabaseStorage) SetValueJSON(ctx context.Context, m models.Metric) int {
	_, err := ds.Database.Exec(ctx, insertValue, m.ID, m.MType, m.Value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, value=%v", err, m.ID, m.MType, m.Value)
		}
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// SetListJSON sets list of metrics in a database by metric list of JSON values
func (ds *DatabaseStorage) SetListJSON(ctx context.Context, list []models.Metric) int {
	tx, err := ds.Database.Begin(ctx)
	if err != nil {
		log.Printf("Error while begin transaction: %v", err)
	}
	for _, v := range list {
		if v.MType == constants.Counter {
			_, err = tx.Exec(ctx, insertDelta, v.ID, v.MType, v.Delta)
		} else {
			_, err = tx.Exec(ctx, insertValue, v.ID, v.MType, v.Value)
		}
		if err != nil {
			log.Printf("%v", err)
			err = tx.Rollback(ctx)
			if err != nil {
				log.Printf("Error while rollback transaction: %v", err)
			}
			return http.StatusInternalServerError
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Error while commit transaction: %v", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// GetAllMetrics returns all metrics from database
func (ds *DatabaseStorage) GetAllMetrics(ctx context.Context) *models.MemStorage {
	result := models.NewMemStorage()
	list := []models.Metric{}
	row := models.Metric{}
	rows, _ := ds.Database.Query(ctx, getAll)
	if rows.Err() != nil {
		log.Printf("Error while execute query: %v", rows.Err())
		return nil
	}
	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&row.ID, &row.MType, &row.Delta, &row.Value)
		if err != nil {
			return nil
		}
		list = append(list, row)
	}

	result.SetListJSON(ctx, list)
	return result
}
