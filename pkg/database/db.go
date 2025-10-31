package database

import (
	"database/sql"
	"time"
)

type Database struct {
	db *sql.DB
}

type AttackLog struct {
	ID         int64
	Timestamp  time.Time
	Service    string
	RemoteAddr string
	Event      string
	Username   string
	Password   string
	Command    string
	RawData    string
}

// NewDatabase initializes a new Database connection.
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	database := &Database{db: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) createTables() error {
	schema := `
		CREATE TABLE IF NOT EXISTS attacks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			service TEXT NOT NULL,
			remote_addr TEXT NOT NULL,
			event TEXT NOT NULL,
			username TEXT,
			password TEXT,
			command TEXT,
			raw_data TEXT
		);

		CREATE INDEX IF NOT EXISTS idx_timestamp ON attacks (timestamp);
		CREATE INDEX IF NOT EXISTS idx_service ON attacks (service);
		CREATE INDEX IF NOT EXISTS idx_remote_addr ON attacks (remote_addr);
	`
	_, err := d.db.Exec(schema)
	return err
}

func (d *Database) InsertAttack(log AttackLog) error {
	query := `
		INSERT INTO attacks (timestamp, service, remote_addr, event, username, password, command, raw_data)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, log.Timestamp, log.Service, log.RemoteAddr, log.Event,
		log.Username, log.Password, log.Command, log.RawData)
	return err
}

// GetRecentAttacks retrieves the most recent attack logs up to the specified limit.
func (d *Database) GetRecentAttacks(limit int) ([]AttackLog, error) {
	query := `
		SELECT id, timestamp, service, remote_addr, event, username, password, command, raw_data
		FROM attacks ORDER BY timestamp DESC LIMIT ?
	`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attacks []AttackLog
	for rows.Next() {
		var attack AttackLog
		err := rows.Scan(&attack.ID, &attack.Timestamp, &attack.Service, &attack.RemoteAddr,
			&attack.Event, &attack.Username, &attack.Password, &attack.Command, &attack.RawData)
		if err != nil {
			return nil, err
		}
		attacks = append(attacks, attack)
	}

	return attacks, nil
}

// GetAttackStats
func (d *Database) GetAttackStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total attacks
	var total int
	d.db.QueryRow("SELECT COUNT(*) FROM attacks").Scan(&total)
	stats["total"] = total

	rows, _ := d.db.Query(`SELECT service, COUNT(*) as count FROM attacks GROUP BY service`)
	defer rows.Close()

	// Attacks by service
	serviceStats := make(map[string]int)
	for rows.Next() {
		var service string
		var count int
		rows.Scan(&service, &count)
		serviceStats[service] = count
	}
	stats["by_service"] = serviceStats

	// Top attacking IPs
	ipRows, _ := d.db.Query(`SELECT remote_addr, COUNT(*) as count FROM attacks GROUP BY remote_addr ORDER BY count DESC LIMIT 10`)
	defer ipRows.Close()

	topIPs := make(map[string]int)
	for ipRows.Next() {
		var ip string
		var count int
		ipRows.Scan(&ip, &count)
		topIPs[ip] = count
	}
	stats["top_ips"] = topIPs

	// Top credentials
	credRows, _ := d.db.Query(`SELECT username, password, COUNT(*) as count FROM attacks WHERE username !='' GROUP BY username, password ORDER BY count DESC LIMIT 10`)
	defer credRows.Close()

	type Credential struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Count    int    `json:"count"`
	}

	// Collect top credentials
	topCreds := []Credential{}
	for credRows.Next() {
		var cred Credential
		credRows.Scan(&cred.Username, &cred.Password, &cred.Count)
		topCreds = append(topCreds, cred)
	}
	stats["top_credentials"] = topCreds

	return stats, nil
}

func (d *Database) Close() {
	d.db.Close()
}
