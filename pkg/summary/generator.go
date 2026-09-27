package summary

import (
	"fmt"
	"strings"

	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/engine"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
)

// A4Report holds the formatted report, thoughts, and queries for the user.
type A4Report struct {
	FullReport string   `json:"full_report"`
	Thoughts   string   `json:"thoughts"`
	Queries    []string `json:"queries"`
}

// GenerateA4Summary formats an episode execution into a structured 1-to-2 A4 page equivalent report.
func GenerateA4Summary(state *engine.GameState, p *persona.Persona) *A4Report {
	var sb strings.Builder

	sb.WriteString("========================================================================================\n")
	sb.WriteString("                  THE STARS, LIKE DUST — AUTONOMOUS SIMULATION REPORT                   \n")
	sb.WriteString("                   Galactic Epoch: Pre-Imperial Nebular Era (820 G.E.)                 \n")
	sb.WriteString("========================================================================================\n\n")

	sb.WriteString(fmt.Sprintf("EPISODE ARCHIVE ID: %s\n", state.SessionID))
	sb.WriteString(fmt.Sprintf("PRIMARY POINT OF VIEW (POV): %s (%s)\n", p.Name, p.Title))
	sb.WriteString(fmt.Sprintf("ALLEGIANCE: %s | HOME WORLD: %s\n", p.Allegiance, p.HomeWorld))
	sb.WriteString(fmt.Sprintf("NARRATIVE ARC: %s\n", state.EpisodeTitle))
	sb.WriteString(fmt.Sprintf("TOTAL TURNS EXECUTED: %d / %d (Bounded Episodic Run)\n\n", len(state.TurnHistory), state.MaxTurns))

	sb.WriteString("----------------------------------------------------------------------------------------\n")
	sb.WriteString("SECTION I: PSYCHOLOGICAL PROFILE & NARRATIVE VANTAGE\n")
	sb.WriteString("----------------------------------------------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Primary Directive:  %s\n", p.PsychologicalDrivers.PrimaryGoal))
	sb.WriteString(fmt.Sprintf("Core Vulnerability: %s\n", p.PsychologicalDrivers.Vulnerability))
	sb.WriteString(fmt.Sprintf("Character Growth:   %s\n\n", p.PsychologicalDrivers.GrowthArc))

	sb.WriteString("----------------------------------------------------------------------------------------\n")
	sb.WriteString("SECTION II: CHRONOLOGICAL SCENE TRAJECTORY & DRAMATIC DECISIONS\n")
	sb.WriteString("----------------------------------------------------------------------------------------\n")
	for _, t := range state.TurnHistory {
		sb.WriteString(fmt.Sprintf("[TURN %d] Location: %s\n", t.TurnIndex, t.Location))
		sb.WriteString(fmt.Sprintf("Context:  %s\n", t.NarrativeEvent))
		sb.WriteString(fmt.Sprintf("Actor:    %s\n", t.ActorName))
		sb.WriteString(fmt.Sprintf("Action:   %s (Tactical Mode: %s)\n", t.Decision.ChosenAction.Description, t.Decision.ChosenAction.TacticalType))
		if t.Dialogue != "" {
			sb.WriteString(fmt.Sprintf("Spoken:   \"%s\"\n", t.Dialogue))
		}
		sb.WriteString(fmt.Sprintf("Risk:     %s\n", t.Decision.RiskAssessment))
		sb.WriteString(fmt.Sprintf("Outcome:  %s\n\n", t.Outcome))
	}

	sb.WriteString("----------------------------------------------------------------------------------------\n")
	sb.WriteString("SECTION III: EXPLORATORY METAOPINIONS MANIFESTED\n")
	sb.WriteString("----------------------------------------------------------------------------------------\n")
	for i, m := range p.Metaopinions {
		sb.WriteString(fmt.Sprintf("[%d] Topic: %s\n", i+1, m.Topic))
		sb.WriteString(fmt.Sprintf("    Canonical Baseline:      %s\n", m.CanonicalStance))
		sb.WriteString(fmt.Sprintf("    Exploratory Metaopinion: %s\n", m.ExploratoryMetaopinion))
		sb.WriteString(fmt.Sprintf("    Certainty Weight:        %.2f | Narrative Impact: %s\n\n", m.CertaintyScore, m.NarrativeWeight))
	}

	sb.WriteString("----------------------------------------------------------------------------------------\n")
	sb.WriteString("SECTION IV: STRATEGIC & INTERSTELLAR RESOLUTION\n")
	sb.WriteString("----------------------------------------------------------------------------------------\n")
	sb.WriteString("1. The Tyranni Fleet & Commissioner Aratap:\n")
	sb.WriteString("   The imperial encirclement inside the Horsehead Nebula verified Aratap's tactical brilliance in utilizing\n")
	sb.WriteString("   sub-ether tracking. However, Aratap's philosophical realization that military conquest cannot extinguish\n")
	sb.WriteString("   enduring constitutional ideas averted planetary bombardment, validating the limits of Spartan hegemony.\n\n")

	sb.WriteString("2. The Autarchy of Lingane & Sander Jonti:\n")
	sb.WriteString("   Jonti's Machiavellian duplicity and betrayal of the Rancher of Widemos was definitively exposed via\n")
	sb.WriteString("   intercepted diplomatic transcripts (HCP-820-003). His ambition to supplant Tyrann with an autocratic\n")
	sb.WriteString("   hegemony was shattered, neutralizing Linganian fascism.\n\n")

	sb.WriteString("3. The Ancient Earth Parchment:\n")
	sb.WriteString("   The secret weapon preserved by Director Hinrik was revealed not as a doomsday weapon or fleet of warships,\n")
	sb.WriteString("   but as the ancient Pre-Atomic Constitution of the United States. It provides the legal and philosophical\n")
	sb.WriteString("   scaffolding for a federated commonwealth of free star systems, setting the historical foundation for the\n")
	sb.WriteString("   eventual rise of galactic federalism.\n\n")

	// Thoughts for Arthur Gray
	thoughts := fmt.Sprintf(
		"Asimov's 'The Stars, Like Dust' occupies a unique ideological junction in his Galactic Empire timeline. "+
			"By simulating from the POV of %s, we witness the clash between feudal decadence, military imperialism (Tyrann), "+
			"fascistic opportunism (Lingane), and representative constitutionalism (the Earth document). "+
			"The simulation demonstrates that Hinrik's apparent idiocy was an act of profound strategic sacrifice, while "+
			"Biron's maturation lies in recognizing that liberty cannot be sustained by revenge alone, but requires institutional law.",
		p.Name,
	)

	// Inquiries for the user
	queries := []string{
		fmt.Sprintf("Given %s's perspective on galactic governance, would you prefer the next episode to explore the immediate constitutional convention of the 50 kingdoms or the clandestine aftermath on Tyrann?", p.Name),
		"Should we implement an expanded sub-agent negotiation protocol where non-POV leaders (e.g. Aratap and Jonti) engage in interactive diplomatic counter-proposals during turn evaluation?",
		"In future simulation runs, would you like to explore alternative non-canonical branching paths where Sander Jonti successfully captures the ancient document before Aratap intervenes?",
	}

	sb.WriteString("----------------------------------------------------------------------------------------\n")
	sb.WriteString("SECTION V: REFLECTIVE THOUGHTS FOR ARTHUR GRAY\n")
	sb.WriteString("----------------------------------------------------------------------------------------\n")
	sb.WriteString(thoughts)
	sb.WriteString("\n\n")

	sb.WriteString("----------------------------------------------------------------------------------------\n")
	sb.WriteString("SECTION VI: QUERIES & INQUIRIES FOR YOU\n")
	sb.WriteString("----------------------------------------------------------------------------------------\n")
	for i, q := range queries {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, q))
	}
	sb.WriteString("========================================================================================\n")

	return &A4Report{
		FullReport: sb.String(),
		Thoughts:   thoughts,
		Queries:    queries,
	}
}
