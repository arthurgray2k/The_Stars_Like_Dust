package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	// MaxDatabaseSizeBytes defines the strict 500 MB hard ceiling
	MaxDatabaseSizeBytes int64 = 500 * 1024 * 1024
)

// Storage handles embedded SQLite database persistence for the simulation.
type Storage struct {
	db     *sql.DB
	dbPath string
}

// Open initializes the embedded database at dbPath, applying migrations.
func Open(dbPath string) (*Storage, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Optimize for single-process embedded simulation
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA foreign_keys=ON;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed executing pragma %s: %w", p, err)
		}
	}

	s := &Storage{db: db, dbPath: dbPath}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed running migrations: %w", err)
	}

	return s, nil
}

// Close closes the underlying database handle.
func (s *Storage) Close() error {
	return s.db.Close()
}

// DB returns the underlying *sql.DB for advanced queries.
func (s *Storage) DB() *sql.DB {
	return s.db
}

// CheckSize verifies that the database file does not exceed the 500 MB limit.
func (s *Storage) CheckSize() (int64, error) {
	info, err := os.Stat(s.dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	if info.Size() > MaxDatabaseSizeBytes {
		return info.Size(), fmt.Errorf("database file %s size (%d bytes) exceeds 500 MB limit", s.dbPath, info.Size())
	}
	return info.Size(), nil
}

func (s *Storage) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS characters (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		title TEXT,
		allegiance TEXT,
		home_world TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS metaopinions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		character_id TEXT NOT NULL,
		topic TEXT NOT NULL,
		canonical_stance TEXT,
		exploratory_metaopinion TEXT NOT NULL,
		certainty_score REAL DEFAULT 0.8,
		narrative_weight TEXT DEFAULT 'high',
		FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS game_sessions (
		id TEXT PRIMARY KEY,
		pov_character TEXT NOT NULL,
		episode_title TEXT NOT NULL,
		started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME,
		status TEXT DEFAULT 'IN_PROGRESS',
		outcome_summary TEXT,
		FOREIGN KEY (pov_character) REFERENCES characters(id)
	);

	CREATE TABLE IF NOT EXISTS session_turns (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		turn_index INTEGER NOT NULL,
		active_actor TEXT NOT NULL,
		location TEXT NOT NULL,
		action_chosen TEXT NOT NULL,
		dialogue TEXT,
		outcome TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (session_id) REFERENCES game_sessions(id) ON DELETE CASCADE
	);

	-- Ephemeral scratch table for branch simulations; pruned after episode summarization
	CREATE TABLE IF NOT EXISTS temp_decision_evaluations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		turn_index INTEGER NOT NULL,
		candidate_action TEXT NOT NULL,
		heuristic_score REAL NOT NULL,
		risk_assessment TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS narrative_summaries (
		session_id TEXT PRIMARY KEY,
		a4_summary_text TEXT NOT NULL,
		thoughts_for_user TEXT,
		queries_for_user TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (session_id) REFERENCES game_sessions(id) ON DELETE CASCADE
	);
	`
	_, err := s.db.Exec(schema)
	return err
}
