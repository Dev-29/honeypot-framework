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
