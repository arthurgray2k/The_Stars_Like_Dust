package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// CharacterRecord models a character row in the database.
type CharacterRecord struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Title      string    `json:"title"`
	Allegiance string    `json:"allegiance"`
	HomeWorld  string    `json:"home_world"`
	CreatedAt  time.Time `json:"created_at"`
}

// MetaopinionRecord models an exploratory or canonical opinion.
type MetaopinionRecord struct {
	ID                     int64   `json:"id"`
	CharacterID            string  `json:"character_id"`
	Topic                  string  `json:"topic"`
	CanonicalStance        string  `json:"canonical_stance"`
	ExploratoryMetaopinion string  `json:"exploratory_metaopinion"`
	CertaintyScore         float64 `json:"certainty_score"`
	NarrativeWeight        string  `json:"narrative_weight"`
}

// SessionRecord models an episodic gameplay session.
type SessionRecord struct {
	ID             string     `json:"id"`
	POVCharacter   string     `json:"pov_character"`
	EpisodeTitle   string     `json:"episode_title"`
	StartedAt      time.Time  `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	Status         string     `json:"status"`
	OutcomeSummary string     `json:"outcome_summary"`
}

// TurnRecord models a single narrative turn.
type TurnRecord struct {
	ID           int64     `json:"id"`
	SessionID    string    `json:"session_id"`
	TurnIndex    int       `json:"turn_index"`
	ActiveActor  string    `json:"active_actor"`
	Location     string    `json:"location"`
	ActionChosen string    `json:"action_chosen"`
	Dialogue     string    `json:"dialogue"`
	Outcome      string    `json:"outcome"`
	CreatedAt    time.Time `json:"created_at"`
}

// UpsertCharacter inserts or updates a character profile.
func (s *Storage) UpsertCharacter(c CharacterRecord) error {
	query := `
	INSERT INTO characters (id, name, title, allegiance, home_world)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		title = excluded.title,
		allegiance = excluded.allegiance,
		home_world = excluded.home_world;
	`
	_, err := s.db.Exec(query, c.ID, c.Name, c.Title, c.Allegiance, c.HomeWorld)
	return err
}

// InsertMetaopinion records an exploratory opinion for a character.
func (s *Storage) InsertMetaopinion(m MetaopinionRecord) error {
	query := `
	INSERT INTO metaopinions (character_id, topic, canonical_stance, exploratory_metaopinion, certainty_score, narrative_weight)
	VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err := s.db.Exec(query, m.CharacterID, m.Topic, m.CanonicalStance, m.ExploratoryMetaopinion, m.CertaintyScore, m.NarrativeWeight)
	return err
}

// GetMetaopinions retrieves all metaopinions registered for a character.
func (s *Storage) GetMetaopinions(characterID string) ([]MetaopinionRecord, error) {
	query := `
	SELECT id, character_id, topic, canonical_stance, exploratory_metaopinion, certainty_score, narrative_weight
	FROM metaopinions
	WHERE character_id = ?;
	`
	rows, err := s.db.Query(query, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []MetaopinionRecord
	for rows.Next() {
		var r MetaopinionRecord
		var canon sql.NullString
		if err := rows.Scan(&r.ID, &r.CharacterID, &r.Topic, &canon, &r.ExploratoryMetaopinion, &r.CertaintyScore, &r.NarrativeWeight); err != nil {
			return nil, err
		}
		if canon.Valid {
			r.CanonicalStance = canon.String
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

// CreateSession registers a new episodic simulation run.
func (s *Storage) CreateSession(sessionID, povCharacter, episodeTitle string) error {
	query := `
	INSERT INTO game_sessions (id, pov_character, episode_title, status)
	VALUES (?, ?, ?, 'IN_PROGRESS');
	`
	_, err := s.db.Exec(query, sessionID, povCharacter, episodeTitle)
	return err
}

// CompleteSession finalizes a session and records its high-level outcome.
func (s *Storage) CompleteSession(sessionID, outcomeSummary string) error {
	query := `
	UPDATE game_sessions
	SET status = 'COMPLETED', outcome_summary = ?, completed_at = CURRENT_TIMESTAMP
	WHERE id = ?;
	`
	_, err := s.db.Exec(query, outcomeSummary, sessionID)
	return err
}

// RecordTurn logs an executed narrative turn.
func (s *Storage) RecordTurn(turn TurnRecord) error {
	query := `
	INSERT INTO session_turns (session_id, turn_index, active_actor, location, action_chosen, dialogue, outcome)
	VALUES (?, ?, ?, ?, ?, ?, ?);
	`
	_, err := s.db.Exec(query, turn.SessionID, turn.TurnIndex, turn.ActiveActor, turn.Location, turn.ActionChosen, turn.Dialogue, turn.Outcome)
	return err
}

// GetTurns retrieves all turns for a session in order.
func (s *Storage) GetTurns(sessionID string) ([]TurnRecord, error) {
	query := `
	SELECT id, session_id, turn_index, active_actor, location, action_chosen, dialogue, outcome, created_at
	FROM session_turns
	WHERE session_id = ?
	ORDER BY turn_index ASC;
	`
	rows, err := s.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turns []TurnRecord
	for rows.Next() {
		var t TurnRecord
		var dial, outc sql.NullString
		if err := rows.Scan(&t.ID, &t.SessionID, &t.TurnIndex, &t.ActiveActor, &t.Location, &t.ActionChosen, &dial, &outc, &t.CreatedAt); err != nil {
			return nil, err
		}
		if dial.Valid {
			t.Dialogue = dial.String
		}
		if outc.Valid {
			t.Outcome = outc.String
		}
		turns = append(turns, t)
	}
	return turns, rows.Err()
}

// RecordTempEvaluation stores an ephemeral analytical calculation during branching simulation.
func (s *Storage) RecordTempEvaluation(sessionID string, turnIndex int, action string, score float64, risk string) error {
	query := `
	INSERT INTO temp_decision_evaluations (session_id, turn_index, candidate_action, heuristic_score, risk_assessment)
	VALUES (?, ?, ?, ?, ?);
	`
	_, err := s.db.Exec(query, sessionID, turnIndex, action, score, risk)
	return err
}

// ClearTempEvaluations purges ephemeral evaluation rows to keep the database lean and under limit.
func (s *Storage) ClearTempEvaluations(sessionID string) error {
	query := `DELETE FROM temp_decision_evaluations WHERE session_id = ?;`
	_, err := s.db.Exec(query, sessionID)
	return err
}

// SaveNarrativeSummary persists the generated A4 summary and user prompts.
func (s *Storage) SaveNarrativeSummary(sessionID, a4Text, thoughts, queries string) error {
	query := `
	INSERT INTO narrative_summaries (session_id, a4_summary_text, thoughts_for_user, queries_for_user)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(session_id) DO UPDATE SET
		a4_summary_text = excluded.a4_summary_text,
		thoughts_for_user = excluded.thoughts_for_user,
		queries_for_user = excluded.queries_for_user;
	`
	_, err := s.db.Exec(query, sessionID, a4Text, thoughts, queries)
	return err
}

// PruneAndVacuum compacts the database and asserts size boundaries.
func (s *Storage) PruneAndVacuum() error {
	if _, err := s.db.Exec("VACUUM;"); err != nil {
		return fmt.Errorf("vacuum failed: %w", err)
	}
	_, err := s.CheckSize()
	return err
}
