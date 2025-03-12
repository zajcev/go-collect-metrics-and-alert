package storage

const (
	getMetric   = "SELECT * FROM metrics WHERE id = $1 and type = $2;"
	getValue    = "SELECT value FROM metrics WHERE id = $1 and type = $2;"
	getDelta    = "SELECT delta FROM metrics WHERE id = $1 and type = $2;"
	getAll      = "SELECT * FROM metrics;"
	insertValue = "INSERT INTO metrics (id, type, value) VALUES ($1, $2, $3) ON CONFLICT(id) DO UPDATE SET value=EXCLUDED.value;"
	insertDelta = "INSERT INTO metrics (id, type, delta) VALUES ($1, $2, $3) ON CONFLICT(id) DO UPDATE SET delta = metrics.delta + EXCLUDED.delta;"
)
