package storage

const (
	getMetric   = "SELECT * FROM metrics WHERE id = $1 and type = $2;"
	getValue    = "SELECT value FROM metrics WHERE id = $1 and type = $2;"
	getDelta    = "SELECT delta FROM metrics WHERE id = $1 and type = $2;"
	getAll      = "SELECT * FROM metrics;"
	updateDelta = "UPDATE metrics SET delta = $1 WHERE id = $2;"
	updateValue = "UPDATE metrics SET value = $1 WHERE id = $2;"
	insertValue = "INSERT INTO metrics (id, type, value) VALUES ($1, $2, $3) ON CONFLICT(id) DO NOTHING;"
	insertDelta = "INSERT INTO metrics (id, type, delta) VALUES ($1, $2, $3) ON CONFLICT(id) DO NOTHING;"
)
