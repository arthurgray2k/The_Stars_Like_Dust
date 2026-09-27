package engine

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/comm"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/storage"
)

// SimulationOptions defines runtime controls for stance, temperature stochasticity, and seeding.
type SimulationOptions struct {
	POV          string  `json:"pov"`
	Stance       string  `json:"stance"`        // balanced, aggressive, covert, diplomatic, inquisitive
	Temperature  float64 `json:"temperature"`   // 0.0 = deterministic; >0.0 = stochastic sampling
	Seed         int64   `json:"seed"`          // 0 = auto-generate time seed
	SubagentMode string  `json:"subagent_mode"` // "off", "archetype", "random"
	SessionID    string  `json:"session_id"`
}

// Engine orchestrates episodic narrative execution.
type Engine struct {
	Store      *storage.Storage
	PersonaMgr *persona.Manager
	CommMgr    *comm.Manager
}

// New creates an initialized simulation engine.
func New(store *storage.Storage, pMgr *persona.Manager, cMgr *comm.Manager) *Engine {
	return &Engine{
		Store:      store,
		PersonaMgr: pMgr,
		CommMgr:    cMgr,
	}
}

// SyncMetadataToDatabase synchronizes loaded persona JSON profiles and metaopinions to the database.
func (e *Engine) SyncMetadataToDatabase() error {
	personas, err := e.PersonaMgr.LoadAll()
	if err != nil {
		return fmt.Errorf("failed loading personas: %w", err)
	}

	for _, p := range personas {
		rec := storage.CharacterRecord{
			ID:         p.ID,
			Name:       p.Name,
			Title:      p.Title,
			Allegiance: p.Allegiance,
			HomeWorld:  p.HomeWorld,
		}
		if err := e.Store.UpsertCharacter(rec); err != nil {
			return fmt.Errorf("failed upserting character %s: %w", p.ID, err)
		}

		for _, m := range p.Metaopinions {
			mRec := storage.MetaopinionRecord{
				CharacterID:            p.ID,
				Topic:                  m.Topic,
				CanonicalStance:        m.CanonicalStance,
				ExploratoryMetaopinion: m.ExploratoryMetaopinion,
				CertaintyScore:         m.CertaintyScore,
				NarrativeWeight:        m.NarrativeWeight,
			}
			if err := e.Store.InsertMetaopinion(mRec); err != nil {
				return fmt.Errorf("failed inserting metaopinion for %s: %w", p.ID, err)
			}
		}
	}
	return nil
}

// RunEpisode executes an episodic simulation with default balanced deterministic settings.
func (e *Engine) RunEpisode(pov string, sessionID string) (*GameState, error) {
	return e.RunEpisodeWithOptions(SimulationOptions{
		POV:          pov,
		Stance:       "balanced",
		Temperature:  0.0,
		Seed:         1,
		SubagentMode: "off",
		SessionID:    sessionID,
	})
}

// RunEpisodeWithOptions executes simulation applying tactical stance and temperature stochasticity.
func (e *Engine) RunEpisodeWithOptions(opts SimulationOptions) (*GameState, error) {
	if err := e.SyncMetadataToDatabase(); err != nil {
		return nil, err
	}

	pov := opts.POV
	povPersona, ok := e.PersonaMgr.Get(pov)
	if !ok {
		pov = "biron_farrill"
		povPersona, _ = e.PersonaMgr.Get(pov)
	}

	sessionID := opts.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%s-%d", pov, time.Now().UnixNano())
	}

	stance := opts.Stance
	if stance == "" {
		stance = "balanced"
	}

	subagentMode := strings.ToLower(strings.TrimSpace(opts.SubagentMode))
	if subagentMode == "true" || subagentMode == "auto" {
		subagentMode = "archetype"
	}
	if subagentMode == "" {
		subagentMode = "off"
	}

	seed := opts.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rnd := rand.New(rand.NewSource(seed))

	episodeTitle := fmt.Sprintf("The Stars, Like Dust: Vantage of %s [%s Stance]", povPersona.Name, stance)
	if err := e.Store.CreateSession(sessionID, pov, episodeTitle); err != nil {
		return nil, fmt.Errorf("failed creating session: %w", err)
	}

	scenes := BuildEpisodeScenes(pov)
	state := &GameState{
		SessionID:      sessionID,
		POVCharacterID: pov,
		EpisodeTitle:   episodeTitle,
		Stance:         stance,
		Temperature:    opts.Temperature,
		Seed:           seed,
		SubagentMode:   subagentMode,
		MaxTurns:       len(scenes),
		Location:       "Starting Orbit",
		ActiveShip:     "The Remembrance",
		EndingBranch:   "Canonical Proclamation of the Free Federation",
		Characters:     make(map[string]*persona.Persona),
		TurnHistory:    make([]TurnResult, 0, len(scenes)),
	}

	for _, p := range e.PersonaMgr.List() {
		state.Characters[p.ID] = p
	}

	if subagentMode != "off" {
		state.Subagents = InitSubagentCast(subagentMode, seed, state.Characters, pov)
	}

	dispatches, _ := e.CommMgr.LoadAll()
	state.Dispatches = dispatches

	for i, scene := range scenes {
		turnIndex := i + 1
		state.CurrentTurn = turnIndex
		state.Location = scene.Location

		actor, exists := e.PersonaMgr.Get(scene.PrimaryActor)
		if !exists {
			actor = povPersona
		}

		// Dynamic consequence adjustments:
		contextDescription := scene.Description
		if scene.Index == 3 && state.HostileCourtTriggered {
			contextDescription = "Due to the explosive confrontation in the Rhodian court, armed sentries patrol Hangar 9 with alert scanners trained on all avenues of approach."
		}

		// Log candidate actions to temporary decision table
		for _, cand := range scene.Options {
			_ = e.Store.RecordTempEvaluation(sessionID, turnIndex, cand.Description, 50.0, cand.TacticalType)
		}

		// Evaluate choice with stance and temperature
		evalCfg := persona.EvaluationConfig{
			Stance:      stance,
			Temperature: opts.Temperature,
			Rnd:         rnd,
		}
		eval := actor.EvaluateChoiceWithConfig(contextDescription, scene.Options, evalCfg)

		// Check for stateful branching flags
		if eval.ChosenAction.ID == "biron_demand_hinrik" {
			state.HostileCourtTriggered = true
		}
		if eval.ChosenAction.ID == "biron_fire_blasters" || eval.ChosenAction.ID == "jonti_feign_innocence" {
			state.BlastersEngaged = true
		}

		outcomeText := fmt.Sprintf("%s resolved scene '%s' (Choice: %s)", actor.Name, scene.Title, eval.ChosenAction.ID)
		var standoff *SubagentStandoff
		if subagentMode != "off" && len(state.Subagents) > 0 {
			standoff = ExecuteSubagentEncounter(turnIndex, scene, actor, eval.ChosenAction, state.Subagents, state.Characters, e.CommMgr, sessionID)
			if standoff != nil {
				outcomeText = fmt.Sprintf("%s resolved '%s' (Choice: %s) -> Inter-Agent Standoff with %s [%s/%s]: %s",
					actor.Name, scene.Title, eval.ChosenAction.ID, standoff.TargetActorName, standoff.SubagentProfile.Model, standoff.SubagentProfile.Effort, standoff.TacticalStroke)
			}
		}

		// Record chosen turn into permanent database table
		turnRec := storage.TurnRecord{
			SessionID:    sessionID,
			TurnIndex:    turnIndex,
			ActiveActor:  actor.ID,
			Location:     scene.Location,
			ActionChosen: eval.ChosenAction.Description,
			Dialogue:     eval.ChosenAction.DialoguePrompt,
			Outcome:      outcomeText,
		}
		if err := e.Store.RecordTurn(turnRec); err != nil {
			return nil, fmt.Errorf("failed recording turn %d: %w", turnIndex, err)
		}

		_ = e.Store.ClearTempEvaluations(sessionID)

		result := TurnResult{
			TurnIndex:        turnIndex,
			ActorID:          actor.ID,
			ActorName:        actor.Name,
			Location:         scene.Location,
			Decision:         eval,
			NarrativeEvent:   contextDescription,
			Dialogue:         eval.ChosenAction.DialoguePrompt,
			Outcome:          turnRec.Outcome,
			SubagentStandoff: standoff,
		}
		state.TurnHistory = append(state.TurnHistory, result)
	}

	// Final outcome resolution based on dynamic branching
	var finalOutcome string
	if state.BlastersEngaged {
		state.EndingBranch = "Aggressive Naval Clash & Ceasefire Compromise"
		finalOutcome = fmt.Sprintf("Episode resolved after direct starship blaster fire. Although Remembrance engaged Aratap's capital cruiser, Aratap's historical rationality averted planetary obliteration, securing an armed truce.")
	} else if stance == "covert" {
		state.EndingBranch = "Covert Archival Exfiltration & Shadow Federation"
		finalOutcome = fmt.Sprintf("Episode resolved through subterranean stealth. Hinrik's ancient constitutional parchment was secured before imperial censors could trace the archival leak.")
	} else {
		state.EndingBranch = "Canonical Proclamation of the Free Federation"
		finalOutcome = fmt.Sprintf("Episode successfully resolved through the decisive agency of %s across %d narrative turns under %s stance.", povPersona.Name, len(scenes), stance)
	}

	if subagentMode != "off" && len(state.Subagents) > 0 {
		finalOutcome = fmt.Sprintf("%s (Sub-Agent Standoffs active: %d secondary actors managed under %s mode)", finalOutcome, len(state.Subagents), subagentMode)
	}

	if err := e.Store.CompleteSession(sessionID, finalOutcome); err != nil {
		return nil, fmt.Errorf("failed completing session: %w", err)
	}

	if err := e.Store.PruneAndVacuum(); err != nil {
		return nil, fmt.Errorf("failed vacuuming database: %w", err)
	}

	return state, nil
}
