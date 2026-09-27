package persona

import (
	"math/rand"
	"path/filepath"
	"testing"
)

func TestPersonaLoadingAndEvaluation(t *testing.T) {
	metadataDir := filepath.Join("..", "..", "metadata")
	mgr := NewManager(metadataDir)

	_, err := mgr.LoadAll()
	if err != nil {
		t.Fatalf("failed to load personas from metadata/: %v", err)
	}

	expectedIDs := []string{"biron_farrill", "artemisia_hinriad", "gillbret_hinriad", "sander_jonti", "simok_aratap", "hinrik_hinriad"}
	for _, id := range expectedIDs {
		p, exists := mgr.Get(id)
		if !exists {
			t.Errorf("expected persona %s to exist in loaded personas", id)
			continue
		}
		if p.Name == "" {
			t.Errorf("persona %s has empty name", id)
		}
		if len(p.Metaopinions) == 0 {
			t.Errorf("persona %s has no metaopinions defined", id)
		}
	}

	all := mgr.List()
	if len(all) < len(expectedIDs) {
		t.Errorf("expected at least %d personas, got %d", len(expectedIDs), len(all))
	}

	// Test decision evaluation for Biron Farrill
	biron, _ := mgr.Get("biron_farrill")
	candidates := []CandidateAction{
		{
			ID:           "act_surrender",
			Description:  "Surrender to the Tyranni border patrol",
			TacticalType: "DIPLOMATIC",
			Keywords:     []string{"surrender", "compromise", "imprisonment"},
			BaseRisk:     0.9,
		},
		{
			ID:           "act_escape_nebula",
			Description:  "Plot hyper-jump into the Horsehead Nebula to avenge Widemos",
			TacticalType: "AGGRESSIVE",
			Keywords:     []string{"avenge", "widemos", "rebellion", "overthrow"},
			BaseRisk:     0.4,
		},
	}

	eval := biron.EvaluateChoice("Confronting imperial blockade", candidates)
	if eval.ChosenAction.ID != "act_escape_nebula" {
		t.Errorf("expected Biron to choose act_escape_nebula, chose %s", eval.ChosenAction.ID)
	}
	if eval.Score <= 0 {
		t.Errorf("expected positive score, got %f", eval.Score)
	}
	if eval.Rationale == "" {
		t.Errorf("expected non-empty rationale")
	}

	// Empty candidates edge case
	emptyEval := biron.EvaluateChoice("Scene", nil)
	if emptyEval.Score != 0 {
		t.Errorf("expected 0 score on empty candidates")
	}
}

func TestPersonaStancesAndTemperature(t *testing.T) {
	metadataDir := filepath.Join("..", "..", "metadata")
	mgr := NewManager(metadataDir)
	_, _ = mgr.LoadAll()
	biron, _ := mgr.Get("biron_farrill")

	candidates := []CandidateAction{
		{
			ID:           "act_diplomatic",
			Description:  "Negotiate an armistice with Commissioner Aratap",
			TacticalType: "DIPLOMATIC",
			Keywords:     []string{"peace", "law", "negotiation"},
			BaseRisk:     0.2,
		},
		{
			ID:           "act_aggressive",
			Description:  "Open fire with forward blaster cannon",
			TacticalType: "AGGRESSIVE",
			Keywords:     []string{"blaster", "attack", "avenge"},
			BaseRisk:     0.8,
		},
		{
			ID:           "act_covert",
			Description:  "Engage stealth dampening and evade sensor probes",
			TacticalType: "COVERT",
			Keywords:     []string{"stealth", "shadow", "evasion"},
			BaseRisk:     0.2,
		},
	}

	// Stance: Aggressive should strongly boost aggressive action
	cfgAggressive := EvaluationConfig{
		Stance:      "aggressive",
		Temperature: 0.0,
	}
	evalAggressive := biron.EvaluateChoiceWithConfig("Tactical engagement", candidates, cfgAggressive)
	if evalAggressive.ChosenAction.ID != "act_aggressive" {
		t.Errorf("expected aggressive stance to choose act_aggressive, chose %s", evalAggressive.ChosenAction.ID)
	}

	// Stance: Covert should boost stealth action
	cfgCovert := EvaluationConfig{
		Stance:      "covert",
		Temperature: 0.0,
	}
	evalCovert := biron.EvaluateChoiceWithConfig("Tactical engagement", candidates, cfgCovert)
	if evalCovert.ChosenAction.ID != "act_covert" {
		t.Errorf("expected covert stance to choose act_covert, chose %s", evalCovert.ChosenAction.ID)
	}

	// Temperature sampling > 0 with fixed seed
	rnd := rand.New(rand.NewSource(42))
	cfgTemp := EvaluationConfig{
		Stance:      "balanced",
		Temperature: 0.8,
		Rnd:         rnd,
	}
	evalTemp := biron.EvaluateChoiceWithConfig("Tactical engagement", candidates, cfgTemp)
	if evalTemp.ChosenAction.ID == "" {
		t.Errorf("expected non-empty sampled action")
	}
	if evalTemp.Probability <= 0 || evalTemp.Probability > 1.0 {
		t.Errorf("expected probability in (0, 1], got %f", evalTemp.Probability)
	}
}
