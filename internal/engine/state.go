package engine

import (
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
	"github.com/arthurgray2k/The_Stars_Like_Dust/pkg/dsl"
)

// GameState encapsulates the simulation runtime state across turns.
type GameState struct {
	SessionID         string                      `json:"session_id"`
	POVCharacterID    string                      `json:"pov_character_id"`
	EpisodeTitle      string                      `json:"episode_title"`
	CurrentTurn       int                         `json:"current_turn"`
	MaxTurns          int                         `json:"max_turns"`
	Location          string                      `json:"location"`
	ActiveShip        string                      `json:"active_ship"`
	TracerAttached    bool                        `json:"tracer_attached"`
	DocumentRecovered bool                        `json:"document_recovered"`
	JontiExposed      bool                        `json:"jonti_exposed"`
	Characters        map[string]*persona.Persona `json:"characters"`
	Dispatches        []*dsl.Dispatch             `json:"dispatches"`
	TurnHistory       []TurnResult                `json:"turn_history"`
}

// TurnResult captures the result of one resolved narrative step.
type TurnResult struct {
	TurnIndex      int                        `json:"turn_index"`
	ActorID        string                     `json:"actor_id"`
	ActorName      string                     `json:"actor_name"`
	Location       string                     `json:"location"`
	Decision       persona.DecisionEvaluation `json:"decision"`
	NarrativeEvent string                     `json:"narrative_event"`
	Dialogue       string                     `json:"dialogue"`
	Outcome        string                     `json:"outcome"`
}
