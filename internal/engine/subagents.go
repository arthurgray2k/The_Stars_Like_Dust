package engine

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/comm"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
	"github.com/arthurgray2k/The_Stars_Like_Dust/pkg/dsl"
)

// InitSubagentCast establishes the subagent delegation profiles for all non-POV actors.
func InitSubagentCast(mode string, seed int64, characters map[string]*persona.Persona, pov string) map[string]SubagentProfile {
	subagents := make(map[string]SubagentProfile)
	if mode == "" || mode == "off" {
		return subagents
	}

	rnd := rand.New(rand.NewSource(seed + 999))
	models := []string{"flash_lite", "flash", "pro"}
	efforts := []string{"fast", "medium", "high"}

	for id, p := range characters {
		if id == pov {
			continue
		}

		profile := SubagentProfile{
			CharacterID: id,
			Name:        p.Name,
		}

		if mode == "random" {
			profile.Model = models[rnd.Intn(len(models))]
			profile.Effort = efforts[rnd.Intn(len(efforts))]
			profile.Role = fmt.Sprintf("Autonomous Sub-Agent (%s)", p.Title)
			profile.Strategy = fmt.Sprintf("Dynamic stochastic response under %s effort", profile.Effort)
		} else {
			// Archetype mode: tailored to Asimovian character role
			switch id {
			case "simok_aratap":
				profile.Model = "pro"
				profile.Effort = "high"
				profile.Role = "Tyranni Chief Inquisitor & Strategic Analyst"
				profile.Strategy = "Cold psychological deduction, fleet surveillance, and sub-ether tracking"
			case "sander_jonti":
				profile.Model = "flash"
				profile.Effort = "medium"
				profile.Role = "Autarch of Lingane & Double Agent"
				profile.Strategy = "Machiavellian ambition, preemptive betrayal, and regional usurpation"
			case "artemisia_hinriad":
				profile.Model = "flash"
				profile.Effort = "high"
				profile.Role = "Rhodian Noble Rebel Diplomat"
				profile.Strategy = "Moral steadfastness, aristocratic alliance, and anti-Tyranni defiance"
			case "hinrik_hinriad":
				profile.Model = "flash_lite"
				profile.Effort = "medium"
				profile.Role = "Director of Rhodia & Archival Guardian"
				profile.Strategy = "Feigned neurotic cowardice concealing generational constitutional preservation"
			case "gillbret_hinriad":
				profile.Model = "flash_lite"
				profile.Effort = "fast"
				profile.Role = "Visisonor Virtuoso & Rebel Astrogator"
				profile.Strategy = "Sensory astrogation, passionate defiance, and lost rebellion lore"
			default:
				profile.Model = "flash"
				profile.Effort = "medium"
				profile.Role = fmt.Sprintf("Interstellar Envoy (%s)", p.Title)
				profile.Strategy = "Pragmatic sector preservation"
			}
		}

		subagents[id] = profile
	}

	return subagents
}

// DetermineTargetActor identifies the primary opposing or interlocutor character in a scene.
func DetermineTargetActor(scene Scene, chosen persona.CandidateAction, pov string) string {
	if chosen.TargetActor != "" {
		return chosen.TargetActor
	}

	// Dynamic contextual inference based on scene context
	switch pov {
	case "biron_farrill":
		switch scene.Index {
		case 1:
			return "sander_jonti" // Jonti secretly sent the communicator warning
		case 2:
			if strings.Contains(chosen.ID, "hinrik") {
				return "hinrik_hinriad"
			}
			return "artemisia_hinriad"
		case 3:
			return "artemisia_hinriad"
		case 4:
			return "sander_jonti"
		case 5:
			if strings.Contains(chosen.ID, "blaster") || strings.Contains(chosen.ID, "constitution") {
				return "simok_aratap"
			}
			return "sander_jonti"
		}
	case "artemisia_hinriad":
		switch scene.Index {
		case 1:
			return "biron_farrill"
		case 2:
			return "hinrik_hinriad"
		case 3:
			return "biron_farrill"
		case 4:
			return "sander_jonti"
		case 5:
			return "simok_aratap"
		}
	case "simok_aratap":
		switch scene.Index {
		case 1, 2, 3:
			return "biron_farrill"
		case 4:
			return "sander_jonti"
		case 5:
			return "biron_farrill"
		}
	case "sander_jonti":
		switch scene.Index {
		case 1, 4:
			return "biron_farrill"
		case 2:
			return "simok_aratap"
		case 3, 5:
			return "simok_aratap"
		}
	case "gillbret_hinriad":
		switch scene.Index {
		case 1, 2, 3:
			return "biron_farrill"
		case 4:
			return "sander_jonti"
		case 5:
			return "simok_aratap"
		}
	case "hinrik_hinriad":
		switch scene.Index {
		case 1, 2:
			return "biron_farrill"
		case 3:
			return "artemisia_hinriad"
		case 4:
			return "simok_aratap"
		case 5:
			return "simok_aratap"
		}
	}

	return "simok_aratap"
}

// ExecuteSubagentEncounter generates the 1-round HCP ping-pong dispatch exchange between protagonist and NPC subagent.
func ExecuteSubagentEncounter(
	turnIndex int,
	scene Scene,
	protagonist *persona.Persona,
	chosenAction persona.CandidateAction,
	subagents map[string]SubagentProfile,
	characters map[string]*persona.Persona,
	commMgr *comm.Manager,
	sessionID string,
) *SubagentStandoff {
	targetID := DetermineTargetActor(scene, chosenAction, protagonist.ID)
	profile, hasSubagent := subagents[targetID]
	if !hasSubagent {
		return nil
	}

	targetPersona := characters[targetID]
	if targetPersona == nil {
		return nil
	}

	// 1. Generate Outbound Protagonist HCP Dispatch
	outboundID := fmt.Sprintf("HCP-%s-T%d", strings.ToUpper(protagonist.ID[:3]), turnIndex)
	outboundDispatch := &dsl.Dispatch{
		ID:         outboundID,
		Encryption: "STANDARD-TACTICAL-CIPHER",
		Origin:     fmt.Sprintf("%s / %s", scene.Location, protagonist.Name),
		Target:     targetPersona.Name,
		Stardate:   "820-G.E.329",
		Priority:   "URGENT",
		Status:     "TRANSMITTED",
		Payload: []dsl.PayloadItem{
			{Type: "EVENT", Content: chosenAction.Description},
			{Type: "INTEL", Content: fmt.Sprintf("Action executed under tactical mode %s.", chosenAction.TacticalType)},
			{Type: "DIRECTIVE", Content: chosenAction.DialoguePrompt},
		},
		Authentication: fmt.Sprintf("SIG-%s-LEAD", strings.ToUpper(strings.ReplaceAll(protagonist.ID, "_", "-"))),
	}

	outboundFilename := fmt.Sprintf("dispatch_%s_t%d_outbound.hcp", sessionID[:8], turnIndex)
	_ = commMgr.SaveDispatch(outboundFilename, outboundDispatch)

	// 2. Synthesize Subagent NPC Counter-Response
	counterID := fmt.Sprintf("HCP-SUB-%s-T%d", strings.ToUpper(targetID[:3]), turnIndex)
	counterEncryption := "COMMON-SECTOR-KEY"
	switch targetID {
	case "simok_aratap":
		counterEncryption = "TYRANN-CIPHER-7"
	case "sander_jonti":
		counterEncryption = "LINGANE-CRYPT-4"
	case "hinrik_hinriad", "artemisia_hinriad":
		counterEncryption = "RHODIA-SEAL-1"
	}

	counterAction, counterDialogue, tacticalStroke := synthesizeSubagentReaction(targetPersona, profile, scene, chosenAction)

	counterDispatch := &dsl.Dispatch{
		ID:         counterID,
		Encryption: counterEncryption,
		Origin:     fmt.Sprintf("%s / %s", targetPersona.Title, targetPersona.Name),
		Target:     protagonist.Name,
		Stardate:   "820-G.E.329",
		Priority:   "FLASH",
		Status:     "INTERCEPTED",
		Headers: map[string]string{
			"SUBAGENT-MODEL":  profile.Model,
			"SUBAGENT-EFFORT": profile.Effort,
		},
		Payload: []dsl.PayloadItem{
			{Type: "EVENT", Content: counterAction},
			{Type: "INTEL", Content: fmt.Sprintf("Sub-agent cognitive evaluation (%s/%s effort) applied.", profile.Model, profile.Effort)},
			{Type: "COUNTER_DEMAND", Content: counterDialogue},
		},
		Authentication: fmt.Sprintf("SIG-SUB-%s-VERIFIED", strings.ToUpper(strings.ReplaceAll(targetID, "_", "-"))),
	}

	counterFilename := fmt.Sprintf("dispatch_%s_t%d_counter.hcp", sessionID[:8], turnIndex)
	_ = commMgr.SaveDispatch(counterFilename, counterDispatch)

	return &SubagentStandoff{
		TurnIndex:       turnIndex,
		ProtagonistID:   protagonist.ID,
		TargetActorID:   targetID,
		TargetActorName: targetPersona.Name,
		SubagentProfile: profile,
		OutboundHCPFile: outboundFilename,
		OutboundAction:  chosenAction.Description,
		CounterHCPFile:  counterFilename,
		CounterAction:   counterAction,
		CounterDialogue: counterDialogue,
		TacticalStroke:  tacticalStroke,
	}
}

func synthesizeSubagentReaction(
	target *persona.Persona,
	profile SubagentProfile,
	scene Scene,
	chosenAction persona.CandidateAction,
) (counterAction, counterDialogue, tacticalStroke string) {
	switch target.ID {
	case "simok_aratap":
		counterAction = "Deploy capital cruiser interdiction web and initiate analytical sub-ether interrogation"
		counterDialogue = "Rebellion cannot triumph against historical necessity. Surrender your coordinates and save your worlds from planetary pacification."
		tacticalStroke = fmt.Sprintf("Aratap's [%s] reasoning isolates protagonist communications, forcing direct terms.", profile.Model)
	case "sander_jonti":
		counterAction = "Assert command over joint expeditionary fleet and feign solidarity against Tyranni patrols"
		counterDialogue = "Lingane alone possesses the cruisers to challenge the Khan. Do not let sentimentality cloud strategic hegemony."
		tacticalStroke = fmt.Sprintf("Jonti's [%s] tactical maneuver attempts to siphon navigational data and neutralize competitors.", profile.Model)
	case "artemisia_hinriad":
		counterAction = "Form dynastic diplomatic coalition and coordinate evasive hyperspace jump"
		counterDialogue = "If House Hinriad must fall, it will fall fighting alongside free men, not groveling as imperial barter."
		tacticalStroke = fmt.Sprintf("Artemisia's [%s] diplomatic resilience stiffens crew morale and secures Rhodian jump beacons.", profile.Model)
	case "hinrik_hinriad":
		counterAction = "Feign nervous panic before Tyranni liaisons while covertly transmitting vault access codes"
		counterDialogue = "Oh dear, the Commissioner's monitors are everywhere! Take the access foil, take it quickly, and speak nothing of this to Tyrann!"
		tacticalStroke = fmt.Sprintf("Hinrik's [%s] deceptive facade distracts imperial censors at the critical juncture.", profile.Model)
	case "gillbret_hinriad":
		counterAction = "Recalibrate visisonor harmonics to mask hyperspace drive signature"
		counterDialogue = "I have seen the nebula in my visions! The stars burn like dust, but there is a free world beyond the dark clouds!"
		tacticalStroke = fmt.Sprintf("Gillbret's [%s] eccentric astrogation masks the ship's ion signature.", profile.Model)
	default:
		counterAction = fmt.Sprintf("Establish standoff containment perimeter under %s directive", profile.Role)
		counterDialogue = "We acknowledge your transmission; comply with sector regulations immediately."
		tacticalStroke = fmt.Sprintf("Sub-agent counter-stroke (%s) deployed.", profile.Model)
	}

	return counterAction, counterDialogue, tacticalStroke
}
