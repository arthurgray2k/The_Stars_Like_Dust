package summary

import (
	"strings"
	"testing"

	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/engine"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
)

func TestGenerateA4Summary(t *testing.T) {
	p := &persona.Persona{
		ID:         "biron_farrill",
		Name:       "Biron Farrill",
		Title:      "Heir of Widemos",
		Allegiance: "Independent",
		HomeWorld:  "Widemos",
		PsychologicalDrivers: persona.PsychologicalDrivers{
			PrimaryGoal:   "Avenge his father and liberate Widemos",
			Vulnerability: "Impulsive trust in allies",
			GrowthArc:     "Embraces constitutional law over mere military vengeance",
		},
		Metaopinions: []persona.Metaopinion{
			{
				Topic:                  "The Nature of the Ultimate Weapon",
				CanonicalStance:        "Expected an armada",
				ExploratoryMetaopinion: "Representative democracy is the true weapon",
				CertaintyScore:         0.92,
				NarrativeWeight:        "pivotal",
			},
		},
	}

	state := &engine.GameState{
		SessionID:      "test-session-001",
		POVCharacterID: "biron_farrill",
		EpisodeTitle:   "Vantage of Biron Farrill",
		Stance:         "aggressive",
		Temperature:    0.7,
		Seed:           999,
		EndingBranch:   "Aggressive Naval Clash & Ceasefire Compromise",
		MaxTurns:       1,
		TurnHistory: []engine.TurnResult{
			{
				TurnIndex: 1,
				ActorID:   "biron_farrill",
				ActorName: "Biron Farrill",
				Location:  "Earth",
				Decision: persona.DecisionEvaluation{
					ChosenAction: persona.CandidateAction{
						Description:  "Evacuate dormitory",
						TacticalType: "COVERT",
					},
					RiskAssessment: "Low operational exposure",
				},
				NarrativeEvent: "Radiation alert triggered",
				Dialogue:       "We must leave Earth now.",
				Outcome:        "Escaped safely.",
			},
		},
	}

	report := GenerateA4Summary(state, p)
	if report == nil {
		t.Fatalf("expected non-nil report")
	}

	if !strings.Contains(report.FullReport, "THE STARS, LIKE DUST — AUTONOMOUS SIMULATION REPORT") {
		t.Errorf("report missing header")
	}
	if !strings.Contains(report.FullReport, "TACTICAL STANCE: AGGRESSIVE") {
		t.Errorf("report missing tactical stance")
	}
	if !strings.Contains(report.FullReport, "TEMPERATURE: 0.70") {
		t.Errorf("report missing temperature")
	}
	if !strings.Contains(report.FullReport, "Aggressive Naval Clash") {
		t.Errorf("report missing branch resolution")
	}
	if !strings.Contains(report.FullReport, "SECTION I: PSYCHOLOGICAL PROFILE") {
		t.Errorf("report missing section I")
	}
	if !strings.Contains(report.FullReport, "SECTION III: EXPLORATORY METAOPINIONS") {
		t.Errorf("report missing section III")
	}
	if !strings.Contains(report.FullReport, "SECTION V: REFLECTIVE THOUGHTS") {
		t.Errorf("report missing section V")
	}
	if !strings.Contains(report.FullReport, "SECTION VI: QUERIES & INQUIRIES") {
		t.Errorf("report missing section VI")
	}
	if len(report.Queries) == 0 {
		t.Errorf("expected queries for user")
	}
	if report.Thoughts == "" {
		t.Errorf("expected thoughts for user")
	}
}
