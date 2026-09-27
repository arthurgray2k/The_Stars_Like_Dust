package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorageLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "stars_db_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "data", "test_universe.db")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer s.Close()

	// Verify size limits
	size, err := s.CheckSize()
	if err != nil {
		t.Fatalf("check size failed: %v", err)
	}
	if size > MaxDatabaseSizeBytes {
		t.Errorf("initial db size %d exceeds max limit %d", size, MaxDatabaseSizeBytes)
	}

	// Test character upsert
	c := CharacterRecord{
		ID:         "biron_farrill",
		Name:       "Biron Farrill",
		Title:      "Heir of Widemos",
		Allegiance: "Independent",
		HomeWorld:  "Widemos",
	}
	if err := s.UpsertCharacter(c); err != nil {
		t.Fatalf("upsert character failed: %v", err)
	}

	// Test metaopinion insert & retrieval
	m := MetaopinionRecord{
		CharacterID:            "biron_farrill",
		Topic:                  "The Nature of the Ultimate Weapon",
		CanonicalStance:        "Expected an armada",
		ExploratoryMetaopinion: "Constitutional democracy is the true weapon",
		CertaintyScore:         0.95,
		NarrativeWeight:        "pivotal",
	}
	if err := s.InsertMetaopinion(m); err != nil {
		t.Fatalf("insert metaopinion failed: %v", err)
	}

	opinions, err := s.GetMetaopinions("biron_farrill")
	if err != nil {
		t.Fatalf("get metaopinions failed: %v", err)
	}
	if len(opinions) != 1 {
		t.Fatalf("expected 1 metaopinion, got %d", len(opinions))
	}
	if opinions[0].Topic != m.Topic {
		t.Errorf("expected topic %s, got %s", m.Topic, opinions[0].Topic)
	}

	// Test game session creation & completion
	sessionID := "sess-test-001"
	if err := s.CreateSession(sessionID, "biron_farrill", "Escape from Rhodia"); err != nil {
		t.Fatalf("create session failed: %v", err)
	}

	// Test recording turns
	turn := TurnRecord{
		SessionID:    sessionID,
		TurnIndex:    1,
		ActiveActor:  "biron_farrill",
		Location:     "Rhodia Directorate Palace",
		ActionChosen: "Confront Director Hinrik",
		Dialogue:     "Where is the document my father spoke of?",
		Outcome:      "Hinrik feigns ignorance while slipping coordinates secretly.",
	}
	if err := s.RecordTurn(turn); err != nil {
		t.Fatalf("record turn failed: %v", err)
	}

	turns, err := s.GetTurns(sessionID)
	if err != nil {
		t.Fatalf("get turns failed: %v", err)
	}
	if len(turns) != 1 {
		t.Fatalf("expected 1 turn, got %d", len(turns))
	}
	if turns[0].ActionChosen != turn.ActionChosen {
		t.Errorf("expected action %s, got %s", turn.ActionChosen, turns[0].ActionChosen)
	}

	// Test temp decision evaluations & purging
	if err := s.RecordTempEvaluation(sessionID, 1, "Flee Palace", 0.72, "High chance of orbital interception"); err != nil {
		t.Fatalf("record temp eval failed: %v", err)
	}
	if err := s.RecordTempEvaluation(sessionID, 1, "Confront Hinrik", 0.89, "Low risk of immediate alarm"); err != nil {
		t.Fatalf("record temp eval failed: %v", err)
	}
	if err := s.ClearTempEvaluations(sessionID); err != nil {
		t.Fatalf("clear temp evaluations failed: %v", err)
	}

	// Complete session
	if err := s.CompleteSession(sessionID, "Escaped successfully into hyperspace."); err != nil {
		t.Fatalf("complete session failed: %v", err)
	}

	// Test narrative summary persistence
	a4Content := "=== A4 EPISODE REPORT ===\nNarrative completed with honors."
	thoughts := "Hinrik's feigned folly is a poignant study in asymmetric survival."
	queries := "Would the rebellion sustain itself if Jonti's betrayal remained hidden?"
	if err := s.SaveNarrativeSummary(sessionID, a4Content, thoughts, queries); err != nil {
		t.Fatalf("save narrative summary failed: %v", err)
	}

	// Test pruning and vacuum
	if err := s.PruneAndVacuum(); err != nil {
		t.Fatalf("prune and vacuum failed: %v", err)
	}

	finalSize, err := s.CheckSize()
	if err != nil {
		t.Fatalf("final size check failed: %v", err)
	}
	if finalSize > MaxDatabaseSizeBytes {
		t.Errorf("database exceeded 500 MB: got %d bytes", finalSize)
	}
}
