package leaderboard

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

const (
	dbDir  = ".visweeper"
	dbFile = "leaderboard.db"
)

type Entry struct {
	ID        int
	Name      string
	Time      time.Duration
	Timestamp time.Time
	Rank      int
}

type Leaderboard struct {
	db *sql.DB
}

func getDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, dbDir, dbFile), nil
}

func New() (*Leaderboard, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Leaderboard{db: db}, nil
}

func (l *Leaderboard) Close() error {
	return l.db.Close()
}

func (l *Leaderboard) createTable(size string) error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		time INTEGER NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`, size)
	_, err := l.db.Exec(query)
	return err
}

func (l *Leaderboard) InsertEntry(size, name string, time time.Duration) error {
	if err := l.createTable(size); err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (name, time) VALUES (?, ?)", size)
	_, err := l.db.Exec(query, name, int64(time.Seconds()))
	return err
}

func (l *Leaderboard) GetAll(table string) ([]Entry, error) {
	if err := l.createTable(table); err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT id, name, time, timestamp FROM %s ORDER BY time ASC LIMIT 1000", table)
	rows, err := l.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	rank := 1
	for rows.Next() {
		var entry Entry
		var timeSeconds int64
		if err := rows.Scan(&entry.ID, &entry.Name, &timeSeconds, &entry.Timestamp); err != nil {
			return nil, err
		}
		entry.Time = time.Duration(timeSeconds) * time.Second
		entry.Rank = rank
		entries = append(entries, entry)
		rank++
	}

	return entries, nil
}

func (l *Leaderboard) Clear(table string) error {
	if err := l.createTable(table); err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s", table)
	_, err := l.db.Exec(query)
	return err
}
