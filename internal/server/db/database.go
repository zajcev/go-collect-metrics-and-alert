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

var db *pgxpool.Pool

//go:embed scripts/000001_test.up.sql
var file string

// Init initializes database connection
func Init(ctx context.Context, DBUrl string) {
	d, err := pgxpool.New(ctx, DBUrl)
	if err != nil {
		log.Printf("Error while connect to database: %v", err)
	}
	db = d
	migration(DBUrl)
}

// migration executes database migration
func migration(DBUrl string) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "internal/server/db/scripts/")
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
func Ping(ctx context.Context) error {
	return db.Ping(ctx)
}

// GetMetricRaw returns metric value from database
func GetMetricRaw(ctx context.Context, mname string, mtype string) interface{} {
	var value interface{}
	if mtype == constants.Gauge {
		row, _ := db.Query(ctx, getValue, mname, mtype)
		if row.Err() != nil {
			log.Printf("Error while execute query: %v", row.Err())
			return nil
		}
		defer row.Close()
		err := row.Scan(&value)
		if err != nil {
			return nil
		}
		return value
	} else {
		row, _ := db.Query(ctx, getDelta, mname, mtype)
		if row.Err() != nil {
			log.Printf("Error while execute query: %v", row.Err())
			return nil
		}
		defer row.Close()
		if row.Next() {
			err := row.Scan(&value)
			if err != nil {
				return nil
			}
			return value
		}
	}
	return nil
}

// GetMetricJSON returns models.Metric for response in JSON format
func GetMetricJSON(ctx context.Context, m models.Metric) (models.Metric, int) {
	row, _ := db.Query(ctx, getMetric, m.ID, m.MType)
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
func SetDeltaRaw(ctx context.Context, mname string, mtype string, delta int64) {
	_, err := db.Exec(ctx, insertDelta, mname, mtype, delta)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, delta=%v", err, mname, mtype, delta)
		}
	}
}

// SetValueRaw sets value for metric type Counter in a database by raw value
func SetValueRaw(ctx context.Context, mname string, mtype string, value float64) {
	_, err := db.Exec(ctx, insertValue, mname, mtype, value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, value=%v", err, mname, mtype, value)
		}
	}
}

// SetDeltaJSON sets value for metric type Gauge in a database by JSON value
func SetDeltaJSON(ctx context.Context, m models.Metric) {
	_, err := db.Exec(ctx, insertDelta, m.ID, m.MType, m.Delta)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, delta=%v", err, m.ID, m.MType, m.Delta)
		}
	}
}

// SetValueJSON sets value for metric type Counter in a database by JSON value
func SetValueJSON(ctx context.Context, m models.Metric) {
	_, err := db.Exec(ctx, insertValue, m.ID, m.MType, m.Value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, value=%v", err, m.ID, m.MType, m.Value)
		}
	}
}

// SetListJSON sets list of metrics in a database by metric list of JSON values
func SetListJSON(ctx context.Context, list []models.Metric) {
	tx, err := db.Begin(ctx)
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
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return
	}
}

// GetAllMetrics returns all metrics from database
func GetAllMetrics(ctx context.Context, ms *models.MemStorage) {
	list := []models.Metric{}
	row := models.Metric{}
	rows, _ := db.Query(ctx, getAll)
	if rows.Err() != nil {
		log.Printf("Error while execute query: %v", rows.Err())
		return
	}
	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&row.ID, &row.MType, &row.Delta, &row.Value)
		if err != nil {
			return
		}
		list = append(list, row)
	}
	ms.SetListJSON(list)
}
