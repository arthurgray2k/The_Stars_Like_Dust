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

func TestEngineSubagentOrchestration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "stars_subagents_test_*")
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

	commDir := filepath.Join(tempDir, "comm")
	cMgr := comm.NewManager(commDir)

	eng := New(store, pMgr, cMgr)

	// Test 1: Archetype Mode
	stateArchetype, err := eng.RunEpisodeWithOptions(SimulationOptions{
		POV:          "biron_farrill",
		Stance:       "balanced",
		Temperature:  0.0,
		Seed:         42,
		SubagentMode: "archetype",
	})
	if err != nil {
		t.Fatalf("run archetype subagents failed: %v", err)
	}

	if len(stateArchetype.Subagents) == 0 {
		t.Fatalf("expected non-empty subagents cast in archetype mode")
	}

	// Verify Aratap archetype profile
	aratapProf, ok := stateArchetype.Subagents["simok_aratap"]
	if !ok {
		t.Fatalf("missing subagent profile for simok_aratap")
	}
	if aratapProf.Model != "pro" || aratapProf.Effort != "high" {
		t.Errorf("expected Aratap to have model pro / high effort, got %s / %s", aratapProf.Model, aratapProf.Effort)
	}

	// Verify Jonti archetype profile
	jontiProf, ok := stateArchetype.Subagents["sander_jonti"]
	if !ok {
		t.Fatalf("missing subagent profile for sander_jonti")
	}
	if jontiProf.Model != "flash" {
		t.Errorf("expected Jonti to have model flash, got %s", jontiProf.Model)
	}

	// Verify standoffs in turns
	hasStandoff := false
	for _, turn := range stateArchetype.TurnHistory {
		if turn.SubagentStandoff != nil {
			hasStandoff = true
			if turn.SubagentStandoff.OutboundHCPFile == "" || turn.SubagentStandoff.CounterHCPFile == "" {
				t.Errorf("turn %d standoff missing HCP filenames", turn.TurnIndex)
			}
			// Verify file was written to disk in commDir
			outPath := filepath.Join(commDir, turn.SubagentStandoff.OutboundHCPFile)
			if _, err := os.Stat(outPath); os.IsNotExist(err) {
				t.Errorf("expected outbound HCP file %s to exist on disk", outPath)
			}
			counterPath := filepath.Join(commDir, turn.SubagentStandoff.CounterHCPFile)
			if _, err := os.Stat(counterPath); os.IsNotExist(err) {
				t.Errorf("expected counter HCP file %s to exist on disk", counterPath)
			}
		}
	}
	if !hasStandoff {
		t.Errorf("expected at least one subagent standoff in episode turns")
	}

	// Test 2: Random Mode
	stateRandom, err := eng.RunEpisodeWithOptions(SimulationOptions{
		POV:          "biron_farrill",
		Stance:       "balanced",
		Temperature:  0.0,
		Seed:         9999,
		SubagentMode: "random",
	})
	if err != nil {
		t.Fatalf("run random subagents failed: %v", err)
	}
	if len(stateRandom.Subagents) == 0 {
		t.Fatalf("expected non-empty subagents cast in random mode")
	}
	for id, prof := range stateRandom.Subagents {
		if prof.Model != "flash_lite" && prof.Model != "flash" && prof.Model != "pro" {
			t.Errorf("subagent %s has invalid random model %s", id, prof.Model)
		}
	}

	// Test 3: Off Mode and Empty Mode
	offCast := InitSubagentCast("off", 1, stateArchetype.Characters, "biron_farrill")
	if len(offCast) != 0 {
		t.Errorf("expected empty cast for off mode")
	}

	// Test 4: Target Determination and Reaction Synthesis across POVs
	testTargets := []struct {
		pov   string
		scene Scene
		opt   persona.CandidateAction
	}{
		{"biron_farrill", Scene{Index: 2}, persona.CandidateAction{ID: "biron_demand_hinrik"}},
		{"artemisia_hinriad", Scene{Index: 1}, persona.CandidateAction{}},
		{"artemisia_hinriad", Scene{Index: 2}, persona.CandidateAction{}},
		{"artemisia_hinriad", Scene{Index: 4}, persona.CandidateAction{}},
		{"simok_aratap", Scene{Index: 4}, persona.CandidateAction{}},
		{"sander_jonti", Scene{Index: 1}, persona.CandidateAction{}},
		{"sander_jonti", Scene{Index: 2}, persona.CandidateAction{}},
		{"gillbret_hinriad", Scene{Index: 4}, persona.CandidateAction{}},
		{"hinrik_hinriad", Scene{Index: 3}, persona.CandidateAction{}},
		{"hinrik_hinriad", Scene{Index: 4}, persona.CandidateAction{}},
	}

	for _, tt := range testTargets {
		target := DetermineTargetActor(tt.scene, tt.opt, tt.pov)
		if target == "" {
			t.Errorf("expected non-empty target for pov %s scene %d", tt.pov, tt.scene.Index)
		}
	}

	// Test reaction synthesis for all persona types
	for _, charID := range []string{"simok_aratap", "sander_jonti", "artemisia_hinriad", "hinrik_hinriad", "gillbret_hinriad", "unknown_char"} {
		testP := &persona.Persona{ID: charID, Name: "Test", Title: "Title"}
		prof := SubagentProfile{Model: "flash", Effort: "medium", Role: "Role"}
		sc := Scene{Index: 1, Title: "Test"}
		opt := persona.CandidateAction{ID: "act", Description: "desc"}
		act, dial, stroke := synthesizeSubagentReaction(testP, prof, sc, opt)
		if act == "" || dial == "" || stroke == "" {
			t.Errorf("expected non-empty reaction for %s", charID)
		}
	}
}
