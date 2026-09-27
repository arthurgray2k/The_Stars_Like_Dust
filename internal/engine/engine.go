package engine

import (
	"fmt"
	"time"

	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/comm"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/storage"
)

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

// RunEpisode executes an episodic, bounded simulation from the chosen character POV.
func (e *Engine) RunEpisode(pov string, sessionID string) (*GameState, error) {
	if err := e.SyncMetadataToDatabase(); err != nil {
		return nil, err
	}

	povPersona, ok := e.PersonaMgr.Get(pov)
	if !ok {
		// Fallback to Biron Farrill if unknown POV passed
		pov = "biron_farrill"
		povPersona, _ = e.PersonaMgr.Get(pov)
	}

	if sessionID == "" {
		sessionID = fmt.Sprintf("session-%s-%d", pov, time.Now().Unix())
	}

	episodeTitle := fmt.Sprintf("The Stars, Like Dust: Vantage of %s", povPersona.Name)
	if err := e.Store.CreateSession(sessionID, pov, episodeTitle); err != nil {
		return nil, fmt.Errorf("failed creating session: %w", err)
	}

	scenes := BuildEpisodeScenes(pov)
	state := &GameState{
		SessionID:      sessionID,
		POVCharacterID: pov,
		EpisodeTitle:   episodeTitle,
		MaxTurns:       len(scenes),
		Location:       "Starting Orbit",
		ActiveShip:     "The Remembrance",
		Characters:     make(map[string]*persona.Persona),
		TurnHistory:    make([]TurnResult, 0, len(scenes)),
	}

	for _, p := range e.PersonaMgr.List() {
		state.Characters[p.ID] = p
	}

	// Load dispatches for reference
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

		// Record candidate evaluations into temporary decision table
		for _, cand := range scene.Options {
			_ = e.Store.RecordTempEvaluation(sessionID, turnIndex, cand.Description, 50.0, cand.TacticalType)
		}

		// Evaluate and select optimal action based on persona
		eval := actor.EvaluateChoice(scene.Description, scene.Options)

		// Record chosen turn into permanent database table
		turnRec := storage.TurnRecord{
			SessionID:    sessionID,
			TurnIndex:    turnIndex,
			ActiveActor:  actor.ID,
			Location:     scene.Location,
			ActionChosen: eval.ChosenAction.Description,
			Dialogue:     eval.ChosenAction.DialoguePrompt,
			Outcome:      fmt.Sprintf("%s successfully resolved scene '%s'", actor.Name, scene.Title),
		}
		if err := e.Store.RecordTurn(turnRec); err != nil {
			return nil, fmt.Errorf("failed recording turn %d: %w", turnIndex, err)
		}

		// Purge ephemeral decision evaluations to keep storage lean
		_ = e.Store.ClearTempEvaluations(sessionID)

		result := TurnResult{
			TurnIndex:      turnIndex,
			ActorID:        actor.ID,
			ActorName:      actor.Name,
			Location:       scene.Location,
			Decision:       eval,
			NarrativeEvent: scene.Description,
			Dialogue:       eval.ChosenAction.DialoguePrompt,
			Outcome:        turnRec.Outcome,
		}
		state.TurnHistory = append(state.TurnHistory, result)
	}

	// Finalize game session in database
	finalOutcome := fmt.Sprintf("Episode successfully resolved through the decisive agency of %s across %d narrative turns.", povPersona.Name, len(scenes))
	if err := e.Store.CompleteSession(sessionID, finalOutcome); err != nil {
		return nil, fmt.Errorf("failed completing session: %w", err)
	}

	// Ensure database remains compact and well under 500 MB
	if err := e.Store.PruneAndVacuum(); err != nil {
		return nil, fmt.Errorf("failed vacuuming database: %w", err)
	}

	return state, nil
}
