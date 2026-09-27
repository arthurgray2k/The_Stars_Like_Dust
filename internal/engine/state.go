package engine

import (
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
	"github.com/arthurgray2k/The_Stars_Like_Dust/pkg/dsl"
)

// GameState encapsulates the simulation runtime state across turns.
type GameState struct {
	SessionID             string                      `json:"session_id"`
	POVCharacterID        string                      `json:"pov_character_id"`
	EpisodeTitle          string                      `json:"episode_title"`
	Stance                string                      `json:"stance"`
	Temperature           float64                     `json:"temperature"`
	Seed                  int64                       `json:"seed"`
	CurrentTurn           int                         `json:"current_turn"`
	MaxTurns              int                         `json:"max_turns"`
	Location              string                      `json:"location"`
	ActiveShip            string                      `json:"active_ship"`
	TracerAttached        bool                        `json:"tracer_attached"`
	DocumentRecovered     bool                        `json:"document_recovered"`
	JontiExposed          bool                        `json:"jonti_exposed"`
	HostileCourtTriggered bool                        `json:"hostile_court_triggered"`
	BlastersEngaged       bool                        `json:"blasters_engaged"`
	EndingBranch          string                      `json:"ending_branch"`
	Characters            map[string]*persona.Persona `json:"characters"`
	SubagentMode          string                      `json:"subagent_mode,omitempty"`
	Subagents             map[string]SubagentProfile  `json:"subagents,omitempty"`
	Dispatches            []*dsl.Dispatch             `json:"dispatches"`
	TurnHistory           []TurnResult                `json:"turn_history"`
}

// SubagentProfile specifies the model tier, effort level, and persona role for an NPC subagent.
type SubagentProfile struct {
	CharacterID string `json:"character_id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Model       string `json:"model"`  // flash_lite, flash, pro
	Effort      string `json:"effort"` // fast, medium, high
	Strategy    string `json:"strategy"`
}

// SubagentStandoff records the 1-round inter-agent interaction between protagonist and NPC.
type SubagentStandoff struct {
	TurnIndex       int             `json:"turn_index"`
	ProtagonistID   string          `json:"protagonist_id"`
	TargetActorID   string          `json:"target_actor_id"`
	TargetActorName string          `json:"target_actor_name"`
	SubagentProfile SubagentProfile `json:"subagent_profile"`
	OutboundHCPFile string          `json:"outbound_hcp_file"`
	OutboundAction  string          `json:"outbound_action"`
	CounterHCPFile  string          `json:"counter_hcp_file"`
	CounterAction   string          `json:"counter_action"`
	CounterDialogue string          `json:"counter_dialogue"`
	TacticalStroke  string          `json:"tactical_stroke"`
}

// TurnResult captures the result of one resolved narrative step.
type TurnResult struct {
	TurnIndex        int                        `json:"turn_index"`
	ActorID          string                     `json:"actor_id"`
	ActorName        string                     `json:"actor_name"`
	Location         string                     `json:"location"`
	Decision         persona.DecisionEvaluation `json:"decision"`
	NarrativeEvent   string                     `json:"narrative_event"`
	Dialogue         string                     `json:"dialogue"`
	Outcome          string                     `json:"outcome"`
	SubagentStandoff *SubagentStandoff          `json:"subagent_standoff,omitempty"`
}
