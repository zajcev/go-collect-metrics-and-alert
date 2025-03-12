package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
	"time"
)

var db *pgxpool.Pool

func Init(DBUrl string) {
	d, err := pgxpool.New(context.Background(), DBUrl)
	if err != nil {
		log.Printf("Error while connect to database: %v", err)
	}
	db = d
}

func Migration() {
	for i := 0; i <= 3; i++ {
		delay := 1
		if i == 3 {
			log.Fatal("Migration failed, stopping after 3 attempts")
		}
		_, err := db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS metrics (id varchar NOT NULL, type varchar NOT NULL,delta bigint NULL,value double precision NULL,CONSTRAINT id UNIQUE (id));")
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "08000" {
				log.Printf("Error: %v| Migration failed, retrying (%d/%d)", err, i+1, 3)
				time.Sleep(time.Duration(delay) * time.Second)
				delay += 2
			} else {
				log.Fatalf("Migration failed: %v", err)
			}
		} else {
			return
		}
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
	row, _ := db.Query(context.Background(), getDelta, mname, mtype)
	defer row.Close()
	if row != nil && row.Next() {
		_, err := db.Exec(context.Background(), updateDelta, delta, mname)
		if err != nil {
			log.Printf("%v", err)
		}
	} else {
		_, err := db.Exec(context.Background(), insertDelta, mname, mtype, delta)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				log.Printf("PGError: %v", pgErr)
			} else {
				log.Printf("Error: %v", err)
			}
		}
	}
}

func SetValueRaw(mname string, mtype string, value float64) {
	row, _ := db.Query(context.Background(), getValue, mname, mtype)
	if row.Err() != nil {
		log.Printf("Error while execute query: %v", row.Err())
	}
	defer row.Close()
	if row != nil && row.Next() {
		_, err := db.Exec(context.Background(), updateValue, value, mname)
		if err != nil {
			log.Printf("%v", err)
		}
	} else {
		_, err := db.Exec(context.Background(), insertValue, mname, mtype, value)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				log.Printf("PGError: %v", pgErr)
			} else {
				log.Printf("Error: %v", err)
			}
		}
	}
}

func SetDeltaJSON(m models.Metric) {
	row, _ := db.Query(context.Background(), getDelta, m.ID, m.MType)
	if row.Err() != nil {
		log.Printf("Error while execute query: %v", row.Err())
	}
	if row != nil && row.Next() {
		res := models.Metric{}
		row.Scan(&res.Delta)
		*m.Delta += *res.Delta
		_, err := db.Exec(context.Background(), updateDelta, m.Delta, m.ID)
		if err != nil {
			log.Printf("%v", err)
		}
	} else {
		_, err := db.Exec(context.Background(), insertDelta, m.ID, m.MType, m.Delta)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				log.Printf("PGError: %v", pgErr)
			} else {
				log.Fatalf("Migration failed: %v", err)
			}
		}
	}
	defer row.Close()
}

func SetValueJSON(m models.Metric) {
	row, _ := db.Query(context.Background(), getValue, m.ID, m.MType)
	if row.Err() != nil {
		log.Printf("Error while execute query: %v", row.Err())
	}
	defer row.Close()
	if row != nil && row.Next() {
		_, err := db.Exec(context.Background(), updateValue, m.Value, m.ID)
		if err != nil {
			log.Printf("%v", err)
		}
	} else {
		_, err := db.Exec(context.Background(), insertValue, m.ID, m.MType, m.Value)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				log.Printf("PGError: %v", pgErr)
			} else {
				log.Fatalf("Migration failed: %v", err)
			}
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
			row, _ := db.Query(context.Background(), getDelta, v.ID, v.MType)
			if row.Err() != nil {
				log.Printf("Error while execute query: %v", row.Err())
			}
			if row.Next() {
				res := models.Metric{}
				err = row.Scan(&res.Delta)
				if err != nil {
					log.Printf("Error while row.Scan: %v", row.Err())
				}
				*v.Delta += *res.Delta
				_, err = tx.Exec(context.Background(), updateDelta, v.Delta, v.ID)
			} else {
				_, err = tx.Exec(context.Background(), insertDelta, v.ID, v.MType, v.Delta)
			}
			row.Close()
		} else {
			row, _ := db.Query(context.Background(), getValue, v.ID, v.MType)
			if row.Next() && row != nil {
				_, err = db.Exec(context.Background(), updateValue, v.Value, v.ID)
				if err != nil {
					log.Printf("%v", err)
				}
			} else {
				_, err = tx.Exec(context.Background(), insertValue, v.ID, v.MType, v.Value)
			}
			row.Close()
		}
		if err != nil {
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
