package engine

import (
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
)

// Scene defines a dramatic situation with context, location, actor, and candidate options.
type Scene struct {
	Index          int                       `json:"index"`
	Title          string                    `json:"title"`
	Location       string                    `json:"location"`
	Description    string                    `json:"description"`
	PrimaryActor   string                    `json:"primary_actor"`
	Options        []persona.CandidateAction `json:"options"`
	NarrativeEvent string                    `json:"narrative_event"`
}

// BuildEpisodeScenes generates the narrative sequence based on chosen character POV.
func BuildEpisodeScenes(pov string) []Scene {
	switch pov {
	case "artemisia_hinriad":
		return buildArtemisiaScenes()
	case "gillbret_hinriad":
		return buildGillbretScenes()
	case "sander_jonti":
		return buildJontiScenes()
	case "simok_aratap":
		return buildSimokAratapScenes()
	case "hinrik_hinriad":
		return buildHinrikScenes()
	case "biron_farrill":
		fallthrough
	default:
		return buildBironScenes()
	}
}

func buildBironScenes() []Scene {
	return []Scene{
		{
			Index:        1,
			Title:        "The Dormitory Trap on Earth",
			Location:     "University of Earth, Sector Solar-0",
			Description:  "A concealed radiation capsule ticks silently in Biron Farrill's dormitory ventilation duct. An urgent anonymous warning buzzes on his communicator.",
			PrimaryActor: "biron_farrill",
			Options: []persona.CandidateAction{
				{
					ID:             "biron_evacuate_earth",
					Description:    "Heed warning, bypass local police, and board cargo transport to Rhodia",
					TacticalType:   "COVERT",
					Keywords:       []string{"avenge", "escape", "rebellion", "widemos"},
					BaseRisk:       0.3,
					DialoguePrompt: "My father taught me that survival precedes victory. Earth is compromised; Rhodia holds the answers.",
				},
				{
					ID:             "biron_confront_police",
					Description:    "Report assassination attempt to Earth Security Commissioner",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"police", "law", "inquiry"},
					BaseRisk:       0.8,
					DialoguePrompt: "Someone planted radioactive poison in my quarters. Who authorized this?",
				},
			},
		},
		{
			Index:        2,
			Title:        "Infiltrating the Directorate of Rhodia",
			Location:     "Hinriad Palace Grounds, Rhodia",
			Description:  "Biron arrives on Rhodia. Director Hinrik appears to tremble in panic, while his daughter Artemisia observes Biron with intense scrutiny.",
			PrimaryActor: "biron_farrill",
			Options: []persona.CandidateAction{
				{
					ID:             "biron_alliance_artemisia",
					Description:    "Propose covert pact with Lady Artemisia to resist Tyranni subjugation",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"rebellion", "honor", "alliance", "constitutional"},
					BaseRisk:       0.35,
					DialoguePrompt: "Lady Artemisia, the Tyranni took my father's life and seek to trade your freedom. We must stand together.",
				},
				{
					ID:             "biron_demand_hinrik",
					Description:    "Publicly demand Director Hinrik explain his subservience to Tyrann",
					TacticalType:   "AGGRESSIVE",
					Keywords:       []string{"confront", "cowardice", "demand", "avenge", "widemos"},
					BaseRisk:       0.75,
					DialoguePrompt: "Does the proud House of Hinriad grovel before desert barbarians?",
				},
			},
		},
		{
			Index:        3,
			Title:        "Seizure of The Remembrance",
			Location:     "Rhodia Imperial Spaceport, Hangar 9",
			Description:  "Accompanied by Artemisia and Gillbret, Biron enters the military hangar where the prototype cruiser Remembrance is prepped for flight.",
			PrimaryActor: "biron_farrill",
			Options: []persona.CandidateAction{
				{
					ID:             "biron_pilot_breakout",
					Description:    "Override hangar blast doors and initiate atmospheric breakout into hyper-space",
					TacticalType:   "AGGRESSIVE",
					Keywords:       []string{"breakout", "pilot", "escape", "widemos"},
					BaseRisk:       0.4,
					DialoguePrompt: "Hold on to the acceleration couches! We break orbit now!",
				},
				{
					ID:             "biron_stealth_launch",
					Description:    "Forge civilian clearance codes and execute sub-orbital glide",
					TacticalType:   "COVERT",
					Keywords:       []string{"stealth", "codes", "glide"},
					BaseRisk:       0.5,
					DialoguePrompt: "Clearance transmitter engaged. Keep all non-essential radiation dampened.",
				},
			},
		},
		{
			Index:        4,
			Title:        "Rendezvous with the Autarch at Lingane",
			Location:     "Deep Space Border, Lingane Sector",
			Description:  "Sander Jonti arrives aboard his armed flagship, offering fuel, arms, and alliance, while demanding complete command over the expedition to the Horsehead Nebula.",
			PrimaryActor: "biron_farrill",
			Options: []persona.CandidateAction{
				{
					ID:             "biron_scrutinize_jonti",
					Description:    "Accept logistics but maintain independent navigational command and probe Jonti's motives",
					TacticalType:   "INQUISITIVE",
					Keywords:       []string{"betrayal", "scrutinize", "widemos", "independence"},
					BaseRisk:       0.3,
					DialoguePrompt: "We welcome Lingane's fuel, Autarch, but Remembrance charts her own course.",
				},
				{
					ID:             "biron_surrender_command",
					Description:    "Surrender navigational control fully to Jonti's military tacticians",
					TacticalType:   "DEFENSIVE",
					Keywords:       []string{"surrender", "subservience", "trust"},
					BaseRisk:       0.85,
					DialoguePrompt: "Take the bridge, Autarch. We place our fates in Lingane's hands.",
				},
			},
		},
		{
			Index:        5,
			Title:        "The Revelation of the Ancient Document",
			Location:     "The Heart of the Horsehead Nebula",
			Description:  "Aratap's imperial fleet springs the trap, tracking Remembrance via sub-ether transponder. In the climactic confrontation, Jonti's betrayal of Widemos is exposed, and Director Hinrik's secret vault discloses the ancient Pre-Atomic Earth document: the democratic Constitution.",
			PrimaryActor: "biron_farrill",
			Options: []persona.CandidateAction{
				{
					ID:             "biron_proclaim_constitution",
					Description:    "Unmask Jonti's treachery, challenge Aratap's imperial order, and proclaim the constitutional federation of free worlds",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"constitutional", "liberty", "federation", "truth", "overthrow"},
					BaseRisk:       0.2,
					DialoguePrompt: "The weapon is not a death ray, Commissioner Aratap! It is the constitutional blueprint of government by the consent of the governed—an idea that will outlive the Khanate!",
				},
				{
					ID:             "biron_fire_blasters",
					Description:    "Engage full blaster batteries against Aratap's capital dreadnought",
					TacticalType:   "AGGRESSIVE",
					Keywords:       []string{"fire", "batteries", "avenge", "widemos", "blasters"},
					BaseRisk:       0.85,
					DialoguePrompt: "Open fire on all forward batteries! Avenge Widemos!",
				},
			},
		},
	}
}

func buildArtemisiaScenes() []Scene {
	return []Scene{
		{
			Index:        1,
			Title:        "The Dynastic Ultimatum",
			Location:     "Directorate Chambers, Rhodia",
			Description:  "Artemisia is confronted with an imperial decree ordering her marriage to a decrepit Tyranni warlord to secure Rhodia's appeasement.",
			PrimaryActor: "artemisia_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "arte_refuse_decree",
					Description:    "Defy the marriage decree and activate covert emergency escape contingency",
					TacticalType:   "AGGRESSIVE",
					Keywords:       []string{"defiance", "liberty", "autonomy", "refuse"},
					BaseRisk:       0.3,
					DialoguePrompt: "I am a daughter of Hinriad, not livestock to barter for the Khan's favor.",
				},
				{
					ID:             "arte_submit_court",
					Description:    "Acquiesce to the court ministers and accept the imperial match",
					TacticalType:   "DEFENSIVE",
					Keywords:       []string{"submit", "appeasement", "marriage"},
					BaseRisk:       0.8,
					DialoguePrompt: "If Rhodia demands my sacrifice, I shall submit.",
				},
			},
		},
		{
			Index:        2,
			Title:        "The Widemos Interception",
			Location:     "Palace Private Gardens, Rhodia",
			Description:  "Artemisia intercepts the fugitive Biron Farrill slipping past the palace guards.",
			PrimaryActor: "artemisia_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "arte_enlist_biron",
					Description:    "Enlist Biron's tactical strength and offer access to cruiser Remembrance",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"alliance", "escape", "resistance", "cruiser"},
					BaseRisk:       0.35,
					DialoguePrompt: "You are Farrill of Widemos. You need a starship, and I need a pilot. Shall we negotiate?",
				},
				{
					ID:             "arte_handover_biron",
					Description:    "Summon palace guards to arrest the Widemos fugitive",
					TacticalType:   "DEFENSIVE",
					Keywords:       []string{"guards", "arrest", "betrayal"},
					BaseRisk:       0.85,
					DialoguePrompt: "Guards! An intruder in the inner courtyard!",
				},
			},
		},
		{
			Index:        3,
			Title:        "Astronavigation into the Void",
			Location:     "Bridge of The Remembrance, Hyperspace",
			Description:  "Artemisia calculates complex hyperspace micro-jumps skirting Tyranni sensor pickets around the Horsehead Nebula.",
			PrimaryActor: "artemisia_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "arte_calculate_jump",
					Description:    "Plot high-precision jump through dark dust corridors relying on Gillbret's sensory logs",
					TacticalType:   "INQUISITIVE",
					Keywords:       []string{"navigation", "corridor", "precision", "nebula"},
					BaseRisk:       0.3,
					DialoguePrompt: "Hold tight. We are diving directly through the dust curtain.",
				},
			},
		},
		{
			Index:        4,
			Title:        "Repelling the Autarch's Ambition",
			Location:     "Wardroom, The Remembrance",
			Description:  "Sander Jonti attempts to persuade Artemisia to unite Rhodia and Lingane under his autocratic banner.",
			PrimaryActor: "artemisia_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "arte_reject_jonti",
					Description:    "Denounce Jonti's predatory ambition and reaffirm Rhodian sovereign independence",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"autonomy", "reject", "sovereignty", "democracy"},
					BaseRisk:       0.25,
					DialoguePrompt: "You do not want a partner, Autarch; you want a vassal planet. Rhodia will never bow to Lingane.",
				},
			},
		},
		{
			Index:        5,
			Title:        "The Vindication of Hinrik",
			Location:     "Flagship Council Chamber, Nebula Sector",
			Description:  "Aratap confronts the prisoners. Hinrik's true master plan is revealed, proving his apparent cowardice was a 30-year sacrifice.",
			PrimaryActor: "artemisia_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "arte_embrace_legacy",
					Description:    "Stand beside her father and proclaim Rhodia's leadership in the democratic federation",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"legacy", "constitutional", "liberty", "reverence"},
					BaseRisk:       0.2,
					DialoguePrompt: "Father, forgive my blindness. Today Rhodia gives humanity the gift of free governance.",
				},
			},
		},
	}
}

func buildSimokAratapScenes() []Scene {
	return []Scene{
		{
			Index:        1,
			Title:        "The Widemos Assessment",
			Location:     "Tyrannic Imperial Headquarters, Rhodia",
			Description:  "Commissioner Aratap reviews intelligence reports following the execution of the Rancher of Widemos. He analyzes potential insurrectionary tremors.",
			PrimaryActor: "simok_aratap",
			Options: []persona.CandidateAction{
				{
					ID:             "aratap_deploy_tracer",
					Description:    "Authorize clandestine installation of sub-ether transponder on cruiser Remembrance",
					TacticalType:   "COVERT",
					Keywords:       []string{"tracer", "surveillance", "order", "calculating"},
					BaseRisk:       0.2,
					DialoguePrompt: "Do not execute Farrill's son. Let him run. He will navigate where our scout probes cannot follow.",
				},
			},
		},
		{
			Index:        2,
			Title:        "Interrogating the Director",
			Location:     "Audience Hall, Rhodia",
			Description:  "Aratap questions Hinrik of Rhodia, observing his nervous tics and medical complaints.",
			PrimaryActor: "simok_aratap",
			Options: []persona.CandidateAction{
				{
					ID:             "aratap_feign_satisfaction",
					Description:    "Politely humor Hinrik's stammering while noting subtle irregularities in archival access logs",
					TacticalType:   "INQUISITIVE",
					Keywords:       []string{"patience", "rationality", "surveillance"},
					BaseRisk:       0.15,
					DialoguePrompt: "Your health concerns are noted, Director. Continue your archival reorganizations at your leisure.",
				},
			},
		},
		{
			Index:        3,
			Title:        "Trailing the Nebula Expedition",
			Location:     "Command Deck, Imperial Dreadnought",
			Description:  "The sub-ether tracer signals Remembrance entering the Horsehead Nebula alongside Lingane's fleet.",
			PrimaryActor: "simok_aratap",
			Options: []persona.CandidateAction{
				{
					ID:             "aratap_shadow_fleet",
					Description:    "Maintain silent sensor distance and let conspirators locate coordinates before intercepting",
					TacticalType:   "COVERT",
					Keywords:       []string{"order", "shadow", "calculating"},
					BaseRisk:       0.2,
					DialoguePrompt: "Hold fire. Let them reach the center. When they drop anchor, we close the net.",
				},
			},
		},
		{
			Index:        4,
			Title:        "Exposing the Linganian Gambit",
			Location:     "Nebular Orbit, Rendezvous Point",
			Description:  "Aratap springs the trap, surrounding Remembrance and Jonti's flagship with twelve heavy cruisers.",
			PrimaryActor: "simok_aratap",
			Options: []persona.CandidateAction{
				{
					ID:             "aratap_expose_jonti",
					Description:    "Reveal Jonti's original signed betrayal of the Rancher of Widemos to splinter the rebels",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"betrayal", "intel", "psychological", "order"},
					BaseRisk:       0.1,
					DialoguePrompt: "Autarch Jonti, did you truly believe the Khanate destroys its intelligence receipts? Young Farrill, look at the transmission log.",
				},
			},
		},
		{
			Index:        5,
			Title:        "The Lesson of the Parchment",
			Location:     "Flagship Council Chamber",
			Description:  "The mysterious rebellion weapon is unveiled as ancient Terra's democratic Constitution. Aratap reflects on the philosophical nature of imperial rule.",
			PrimaryActor: "simok_aratap",
			Options: []persona.CandidateAction{
				{
					ID:             "aratap_philosophical_verdict",
					Description:    "Acknowledge the historical inevitability of democratic integration and avoid planetary slaughter",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"rationality", "historical", "governance", "future"},
					BaseRisk:       0.2,
					DialoguePrompt: "A paper constitution... No battleships, no planet-busters. You may keep your ancient parchment, Farrill. Let us see if human nature can live up to it.",
				},
			},
		},
	}
}

func buildGillbretScenes() []Scene {
	return []Scene{
		{
			Index:        1,
			Title:        "The Visisonor Recital",
			Location:     "Rhodian Court Ballroom",
			Description:  "Gillbret performs on his visisonor before the visiting Tyranni military inspectors, weaving hypnotic patterns of light and music.",
			PrimaryActor: "gillbret_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "gill_visisonor_hypnosis",
					Description:    "Project emotional harmonics to induce disorientation among the Tyranni guard officers",
					TacticalType:   "COVERT",
					Keywords:       []string{"visisonor", "music", "psychology", "genius"},
					BaseRisk:       0.25,
					DialoguePrompt: "Listen to the light, gentlemen... feel the harmony of the stars.",
				},
			},
		},
		{
			Index:        2,
			Title:        "The Secret of the Shipwreck",
			Location:     "Gillbret's Private Observatory, Rhodia",
			Description:  "Gillbret reveals the truth of his twenty-year-old hyperspace shipwreck to Biron and Artemisia.",
			PrimaryActor: "gillbret_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "gill_share_coordinates",
					Description:    "Disclose the encoded astrogation logs pointing into the Horsehead Nebula",
					TacticalType:   "INQUISITIVE",
					Keywords:       []string{"coordinates", "nebula", "truth", "rebellion"},
					BaseRisk:       0.3,
					DialoguePrompt: "They called me mad! But I saw it with my own eyes—lights and vessels thriving in the deep dust!",
				},
			},
		},
		{
			Index:        3,
			Title:        "Sabotage on the Hangar Floor",
			Location:     "Rhodia Spaceport Maintenance Deck",
			Description:  "Gillbret uses acoustic tools to disable the port tracking radar array.",
			PrimaryActor: "gillbret_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "gill_disable_radar",
					Description:    "Overload radar receiver with harmonic acoustic feedback",
					TacticalType:   "COVERT",
					Keywords:       []string{"sabotage", "acoustic", "escape"},
					BaseRisk:       0.35,
					DialoguePrompt: "A simple harmonic surge in the waveguide... and Tyrann's eyes go blind.",
				},
			},
		},
		{
			Index:        4,
			Title:        "The Visisonor Duel of Minds",
			Location:     "Nebular Orbit, Aboard The Remembrance",
			Description:  "Gillbret detects deceit in Sander Jonti's aura and uses the visisonor to test the Autarch's subconscious.",
			PrimaryActor: "gillbret_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "gill_probe_jonti",
					Description:    "Channel cognitive probe chords to reveal Jonti's suppressed guilt over Widemos",
					TacticalType:   "INQUISITIVE",
					Keywords:       []string{"probe", "guilt", "visisonor", "exposure"},
					BaseRisk:       0.3,
					DialoguePrompt: "Why does the Autarch flinch at the sound of the Widemos mourning chord?",
				},
			},
		},
		{
			Index:        5,
			Title:        "Vindication in the Dark Cloud",
			Location:     "Flagship Council Chamber",
			Description:  "Standing before Aratap and Hinrik, Gillbret realizes the true meaning of his vision.",
			PrimaryActor: "gillbret_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "gill_claim_vindication",
					Description:    "Celebrate the vindication of his sanity and the awakening of fifty star kingdoms",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"vindication", "sanity", "liberty", "stars"},
					BaseRisk:       0.2,
					DialoguePrompt: "I was never mad! The rebellion world was not in the dust—it was in the hearts of free men!",
				},
			},
		},
	}
}

func buildJontiScenes() []Scene {
	return []Scene{
		{
			Index:        1,
			Title:        "The Elimination of Widemos",
			Location:     "Autarch's Citadel, Lingane",
			Description:  "Sander Jonti reflects on delivering the sealed dossier on the Rancher of Widemos to Tyranni intelligence.",
			PrimaryActor: "sander_jonti",
			Options: []persona.CandidateAction{
				{
					ID:             "jonti_justify_betrayal",
					Description:    "Consolidate control over regional rebel cells now that Widemos is removed",
					TacticalType:   "COVERT",
					Keywords:       []string{"autocracy", "hegemony", "consolidation", "realpolitik"},
					BaseRisk:       0.3,
					DialoguePrompt: "Widemos was a sentimental fool. A galactic revolution requires an iron fist, not a romantic rancher.",
				},
			},
		},
		{
			Index:        2,
			Title:        "Guiding the Widemos Heir",
			Location:     "Subspace Communication Hub, Lingane",
			Description:  "Jonti dispatches covert alerts to Earth to push young Biron toward Rhodia.",
			PrimaryActor: "sander_jonti",
			Options: []persona.CandidateAction{
				{
					ID:             "jonti_manipulate_biron",
					Description:    "Transmit anonymous coordinates to ensure Biron stirs up Rhodian court intrigue",
					TacticalType:   "COVERT",
					Keywords:       []string{"manipulate", "pawn", "rhodia"},
					BaseRisk:       0.25,
					DialoguePrompt: "Let young Farrill play the avenging hero and draw Aratap's gaze while we ready the fleet.",
				},
			},
		},
		{
			Index:        3,
			Title:        "The Demand for Supremacy",
			Location:     "Space Sector Lingane-Outpost",
			Description:  "Jonti boards Remembrance and attempts to coerce Artemisia and Biron into acknowledging his supreme authority.",
			PrimaryActor: "sander_jonti",
			Options: []persona.CandidateAction{
				{
					ID:             "jonti_assert_leadership",
					Description:    "Demand sole command of the Horsehead expedition in exchange for Linganian escort",
					TacticalType:   "AGGRESSIVE",
					Keywords:       []string{"autocracy", "command", "supremacy"},
					BaseRisk:       0.4,
					DialoguePrompt: "I command six battle cruisers. You possess a stolen prototype. Do not mistake courtesy for weakness.",
				},
			},
		},
		{
			Index:        4,
			Title:        "The Unraveling Conspiracy",
			Location:     "Inside the Horsehead Nebula",
			Description:  "Aratap's imperial battlefleet surrounds the flotilla, and the sub-ether tracking device is identified.",
			PrimaryActor: "sander_jonti",
			Options: []persona.CandidateAction{
				{
					ID:             "jonti_feign_innocence",
					Description:    "Blame Rhodia for the sub-ether bug and draw blaster to seize Biron as hostage",
					TacticalType:   "AGGRESSIVE",
					Keywords:       []string{"hostage", "blaster", "survival"},
					BaseRisk:       0.8,
					DialoguePrompt: "Farrill! You led them to us! Back away or I disintegrate the girl!",
				},
			},
		},
		{
			Index:        5,
			Title:        "The Fall of the Autarch",
			Location:     "Imperial Interrogation Deck",
			Description:  "Disarmed and exposed by Aratap's decrypted receipts, Jonti faces total political ruin.",
			PrimaryActor: "sander_jonti",
			Options: []persona.CandidateAction{
				{
					ID:             "jonti_final_curse",
					Description:    "Curse democratic idealism and predict the collapse of all fifty kingdoms into barbarian ruin",
					TacticalType:   "AGGRESSIVE",
					Keywords:       []string{"autocracy", "cynicism", "ruin"},
					BaseRisk:       0.2,
					DialoguePrompt: "You are children! When the Khanate falls, your 'democracy' will drown in blood within a generation!",
				},
			},
		},
	}
}

func buildHinrikScenes() []Scene {
	return []Scene{
		{
			Index:        1,
			Title:        "The Mask of Incompetence",
			Location:     "Director's Solar, Rhodia",
			Description:  "Director Hinrik rehearses his nervous stammer and prepares for Commissioner Aratap's unscheduled audit.",
			PrimaryActor: "hinrik_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "hinrik_play_fool",
					Description:    "Fidget frantically with medical tonics and complaints to divert Aratap from archival activities",
					TacticalType:   "COVERT",
					Keywords:       []string{"mask", "folly", "survival", "camouflage"},
					BaseRisk:       0.1,
					DialoguePrompt: "Oh dear, my nervous palpitations again... where did the physician place the sedative drops?",
				},
			},
		},
		{
			Index:        2,
			Title:        "Shielding Artemisia",
			Location:     "Private Corridor, Rhodian Palace",
			Description:  "Hinrik overhears his daughter Artemisia expressing bitter shame over his apparent cowardice.",
			PrimaryActor: "hinrik_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "hinrik_endure_scorn",
					Description:    "Swallow the agony of his daughter's contempt to keep her off the Tyranni execution list",
					TacticalType:   "DEFENSIVE",
					Keywords:       []string{"sacrifice", "patience", "protection"},
					BaseRisk:       0.15,
					DialoguePrompt: "Despise me if you must, my sweet child... but live. Live until the day the chains break.",
				},
			},
		},
		{
			Index:        3,
			Title:        "Unlocking the Hangar Vault",
			Location:     "Deep Vault beneath Spaceport 9",
			Description:  "Hinrik secretly overrides the security interlocks on prototype cruiser Remembrance, ensuring Biron can access it.",
			PrimaryActor: "hinrik_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "hinrik_enable_escape",
					Description:    "Key in the master clearance bypass and slip the coordinates into Biron's navigation console",
					TacticalType:   "COVERT",
					Keywords:       []string{"bypass", "vault", "rebellion", "constitutional"},
					BaseRisk:       0.25,
					DialoguePrompt: "Widemos died for this dream. Go, young Biron, and carry humanity forward.",
				},
			},
		},
		{
			Index:        4,
			Title:        "The Great Deception Disclosed",
			Location:     "Directorate Archive Chamber",
			Description:  "Aratap demands the surrender of the suspected rebellion base coordinates.",
			PrimaryActor: "hinrik_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "hinrik_deliver_document",
					Description:    "Drop the stammering facade, straighten his spine, and present the ancient Pre-Atomic Earth Constitution",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"constitutional", "revelation", "truth", "sovereignty"},
					BaseRisk:       0.2,
					DialoguePrompt: "Look well upon me, Commissioner Aratap. For thirty years I have played the clown so this parchment might live. Here is the true weapon.",
				},
			},
		},
		{
			Index:        5,
			Title:        "The New Dawn for the Fifty Kingdoms",
			Location:     "Directorate Balcony, Rhodia",
			Description:  "With Jonti disgraced and Aratap withdrawing his garrison, Hinrik addresses the free worlds.",
			PrimaryActor: "hinrik_hinriad",
			Options: []persona.CandidateAction{
				{
					ID:             "hinrik_federation_dawn",
					Description:    "Proclaim the constitutional assembly of the fifty star kingdoms",
					TacticalType:   "DIPLOMATIC",
					Keywords:       []string{"federation", "democracy", "constitution", "dawn"},
					BaseRisk:       0.1,
					DialoguePrompt: "Let the stars shine like dust across the black sky—not as isolated tyrannies, but as a commonwealth of free peoples.",
				},
			},
		},
	}
}
