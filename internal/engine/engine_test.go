package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/comm"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/storage"
)

func TestEngineSimulationAcrossPOVs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "stars_engine_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "data", "test_universe.db")
	store, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer store.Close()

	metadataDir := filepath.Join("..", "..", "metadata")
	pMgr := persona.NewManager(metadataDir)

	commDir := filepath.Join("..", "..", "comm")
	cMgr := comm.NewManager(commDir)

	eng := New(store, pMgr, cMgr)

	testPOVs := []string{
		"biron_farrill",
		"artemisia_hinriad",
		"simok_aratap",
		"gillbret_hinriad",
		"sander_jonti",
		"hinrik_hinriad",
	}

	for _, pov := range testPOVs {
		t.Run("POV_"+pov, func(t *testing.T) {
			state, err := eng.RunEpisode(pov, "")
			if err != nil {
				t.Fatalf("run episode failed for pov %s: %v", pov, err)
			}

			if state.POVCharacterID != pov {
				t.Errorf("expected pov %s, got %s", pov, state.POVCharacterID)
			}
			if len(state.TurnHistory) == 0 {
				t.Fatalf("expected non-empty turn history for %s", pov)
			}
			if state.CurrentTurn != state.MaxTurns {
				t.Errorf("expected current turn %d to equal max turns %d", state.CurrentTurn, state.MaxTurns)
			}

			// Verify turns stored in database
			turns, err := store.GetTurns(state.SessionID)
			if err != nil {
				t.Fatalf("failed to get turns from store: %v", err)
			}
			if len(turns) != len(state.TurnHistory) {
				t.Errorf("expected %d turns in db, got %d", len(state.TurnHistory), len(turns))
			}

			// Verify database file size remains tiny (< 500 MB)
			size, err := store.CheckSize()
			if err != nil {
				t.Fatalf("failed checking db size: %v", err)
			}
			if size > storage.MaxDatabaseSizeBytes {
				t.Errorf("database size %d exceeds limit", size)
			}
		})
	}
}

func TestEngineDivergentOutcomesSameCharacter(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "stars_divergence_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "data", "test_universe.db")
	store, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer store.Close()

	metadataDir := filepath.Join("..", "..", "metadata")
	pMgr := persona.NewManager(metadataDir)

	commDir := filepath.Join("..", "..", "comm")
	cMgr := comm.NewManager(commDir)

	eng := New(store, pMgr, cMgr)

	// Run 1: Balanced stance (canonical)
	stateBalanced, err := eng.RunEpisodeWithOptions(SimulationOptions{
		POV:         "biron_farrill",
		Stance:      "balanced",
		Temperature: 0.0,
		Seed:        100,
	})
	if err != nil {
		t.Fatalf("run balanced failed: %v", err)
	}

	// Run 2: Aggressive stance
	stateAggressive, err := eng.RunEpisodeWithOptions(SimulationOptions{
		POV:         "biron_farrill",
		Stance:      "aggressive",
		Temperature: 0.0,
		Seed:        100,
	})
	if err != nil {
		t.Fatalf("run aggressive failed: %v", err)
	}

	// Run 3: Covert stance
	stateCovert, err := eng.RunEpisodeWithOptions(SimulationOptions{
		POV:         "biron_farrill",
		Stance:      "covert",
		Temperature: 0.0,
		Seed:        100,
	})
	if err != nil {
		t.Fatalf("run covert failed: %v", err)
	}

	// Verify that aggressive stance triggered different choices than balanced
	// For instance, turn 2 (Hinriad court confrontation) or turn 5 (blaster combat)
	if stateBalanced.EndingBranch == stateAggressive.EndingBranch {
		t.Errorf("expected different ending branch between balanced (%s) and aggressive (%s)",
			stateBalanced.EndingBranch, stateAggressive.EndingBranch)
	}

	if stateCovert.EndingBranch == stateAggressive.EndingBranch {
		t.Errorf("expected different ending branch between covert (%s) and aggressive (%s)",
			stateCovert.EndingBranch, stateAggressive.EndingBranch)
	}

	// Verify stochastic run with temperature > 0
	stateStochastic, err := eng.RunEpisodeWithOptions(SimulationOptions{
		POV:         "biron_farrill",
		Stance:      "balanced",
		Temperature: 0.9,
		Seed:        42,
	})
	if err != nil {
		t.Fatalf("run stochastic failed: %v", err)
	}
	if len(stateStochastic.TurnHistory) != 5 {
		t.Errorf("expected 5 turns, got %d", len(stateStochastic.TurnHistory))
	}
}
