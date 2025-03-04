package storage

import (
	"database/sql"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/constants"
	"github.com/zajcev/go-collect-metrics-and-alert/internal/server/models"
	"log"
	"net/http"
)

var db *sql.DB

func Init(DBUrl string) {
	d, err := sql.Open("postgres", DBUrl)
	if err != nil {
		log.Printf("Error while connect to database: %v", err)
	}
	db = d
}

func Migration() {
	db.Exec("CREATE TABLE IF NOT EXISTS metrics (id varchar NOT NULL,mtype varchar NOT NULL,delta bigint NULL,value double precision NULL,CONSTRAINT metrics_pk PRIMARY KEY (id),CONSTRAINT metrics_unique UNIQUE (name));")
}

func DBPing() error {
	return db.Ping()
}

func GetMetricRaw(mname string, mtype string) interface{} {
	var value interface{}
	if mtype == constants.Gauge {
		row, _ := db.Query("SELECT value FROM metrics WHERE id = $1 and type = $2;", mname, mtype)
		defer row.Close()
		if row != nil {
			row.Scan(&value)
			return value
		}
	} else {
		row, _ := db.Query("SELECT value FROM metrics WHERE id = $1 and type = $2;", mname, mtype)
		defer row.Close()
		if row != nil {
			row.Scan(&value)
			return value
		}
	}
	return nil
}

func GetMetricJSON(m models.Metric) (models.Metric, int) {
	row, _ := db.Query("SELECT value FROM metrics WHERE id = $1;", m.ID)
	defer row.Close()
	if row.Next() {
		row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		return m, http.StatusOK
	}
	return models.Metric{}, http.StatusNotFound
}

func SetDeltaRaw(mname string, mtype string, delta int64) {
	row, _ := db.Query("SELECT * FROM metrics WHERE id = $1 and type = $2;", mname, mtype)
	defer row.Close()
	if row.Next() {
		_, err := db.Exec("UPDATE metrics SET delta = $1 WHERE id = $2;", delta, mname)
		if err != nil {
			log.Printf("%v", err)
		}
	} else {
		_, err := db.Exec("INSERT INTO metrics (id, type, delta) VALUES ($1, $2, $3);", mname, mtype, delta)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}

func SetValueRaw(mname string, mtype string, value float64) {
	row, _ := db.Query("SELECT * FROM metrics WHERE id = $1 and type = $2;", mname, mtype)
	defer row.Close()
	if row.Next() {
		_, err := db.Exec("UPDATE metrics SET value = $1 WHERE id = $2;", value, mname)
		if err != nil {
			log.Printf("%v", err)
		}
	} else {
		_, err := db.Exec("INSERT INTO metrics (id, type, value) VALUES ($1, $2, $3);", mname, mtype, value)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}

func SetDeltaJSON(m models.Metric) {
	row, _ := db.Query("SELECT * FROM metrics WHERE id = $1 and type = $2;", m.ID, m.MType)
	if row.Next() {
		_, err := db.Exec("UPDATE metrics SET delta = $1 WHERE id = $2;", m.Delta, m.ID)
		if err != nil {
			log.Printf("%v", err)
		}
	} else {
		_, err := db.Exec("INSERT INTO metrics (id, type, delta) VALUES ($1, $2, $3);", m.ID, m.MType, m.Delta)
		if err != nil {
			log.Printf("Error while insert metric with counter type: %v", err)
		}
	}
	defer row.Close()
}

func SetValueJSON(m models.Metric) {
	row, _ := db.Query("SELECT * FROM metrics WHERE id = $1 and type = $2;", m.ID, m.MType)
	defer row.Close()
	if row.Next() {
		_, err := db.Exec("UPDATE metrics SET delta = $1 WHERE id = $2;", m.Delta, m.ID)
		if err != nil {
			log.Printf("%v", err)
		}
	} else {
		result, err := db.Exec("INSERT INTO metrics (id, type, value) VALUES ($1, $2, $3);", m.ID, m.MType, m.Value)
		log.Print(result)
		if err != nil {
			log.Printf("Error while insert metric with gauge type %v", err)
		}
	}
}

func GetAllMetrics() (*models.MemStorage, error) {
	metric := models.MemStorage{}
	rows, err := db.Query("SELECT * FROM metrics;")
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	defer rows.Close()
	rows.Scan(metric)
	return &metric, nil
}
