package db

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var db *pgxpool.Pool

//go:embed scripts/000001_test.up.sql
var file string

func Init(DBUrl string) {
	d, err := pgxpool.New(context.Background(), DBUrl)
	if err != nil {
		log.Printf("Error while connect to database: %v", err)
	}
	db = d
	migration(DBUrl)
}

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

func DBPing() error {
	return db.Ping(context.Background())
}

func GetMetricRaw(mname string, mtype string) interface{} {
	var value interface{}
	if mtype == constants.Gauge {
		row, _ := db.Query(context.Background(), getValue, mname, mtype)
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
		row, _ := db.Query(context.Background(), getDelta, mname, mtype)
		if row.Err() != nil {
			log.Printf("Error while execute query: %v", row.Err())
			return nil
		}
		defer row.Close()
		if row.Next() {
			row.Scan(&value)
			return value
		}
	}
	return nil
}

func GetMetricJSON(m models.Metric) (models.Metric, int) {
	row, _ := db.Query(context.Background(), getMetric, m.ID, m.MType)
	if row.Err() != nil {
		log.Printf("Error while execute query: %v", row.Err())
	}
	defer row.Close()
	if row.Next() {
		row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		return m, http.StatusOK
	}
	return models.Metric{}, http.StatusNotFound
}

func SetDeltaRaw(mname string, mtype string, delta int64) {
	_, err := db.Exec(context.Background(), insertDelta, mname, mtype, delta)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, delta=%v", err, mname, mtype, delta)
		}
	}
}

func SetValueRaw(mname string, mtype string, value float64) {
	_, err := db.Exec(context.Background(), insertValue, mname, mtype, value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, value=%v", err, mname, mtype, value)
		}
	}
}

func SetDeltaJSON(m models.Metric) {
	_, err := db.Exec(context.Background(), insertDelta, m.ID, m.MType, m.Delta)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, delta=%v", err, m.ID, m.MType, m.Delta)
		}
	}
}

func SetValueJSON(m models.Metric) {
	_, err := db.Exec(context.Background(), insertValue, m.ID, m.MType, m.Value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while insert metric: Error=%v, id=%v, type=%v, value=%v", err, m.ID, m.MType, m.Value)
		}
	}
}

func SetListJSON(list []models.Metric) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Printf("Error while begin transaction: %v", err)
	}
	for _, v := range list {
		if v.MType == constants.Counter {
			_, err = tx.Exec(context.Background(), insertDelta, v.ID, v.MType, v.Delta)
		} else {
			_, err = tx.Exec(context.Background(), insertValue, v.ID, v.MType, v.Value)
		}
		if err != nil {
			log.Printf("%v", err)
			err = tx.Rollback(context.Background())
			if err != nil {
				log.Printf("Error while rollback transaction: %v", err)
			}
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return
	}
}

func GetAllMetrics(ms *models.MemStorage) {
	list := []models.Metric{}
	row := models.Metric{}
	rows, _ := db.Query(context.Background(), getAll)
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
	ms.SetMetricList(list)
}
