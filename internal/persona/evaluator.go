package persona

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

// CandidateAction defines an action option available during a narrative scene.
type CandidateAction struct {
	ID             string   `json:"id"`
	Description    string   `json:"description"`
	TacticalType   string   `json:"tactical_type"` // DIPLOMATIC, COVERT, AGGRESSIVE, DEFENSIVE, INQUISITIVE
	TargetActor    string   `json:"target_actor"`
	Keywords       []string `json:"keywords"`
	BaseRisk       float64  `json:"base_risk"` // 0.0 to 1.0
	DialoguePrompt string   `json:"dialogue_prompt"`
}

// DecisionEvaluation records the scoring and rationale of a chosen action.
type DecisionEvaluation struct {
	ChosenAction   CandidateAction `json:"chosen_action"`
	Score          float64         `json:"score"`
	Probability    float64         `json:"probability,omitempty"`
	RiskAssessment string          `json:"risk_assessment"`
	Rationale      string          `json:"rationale"`
	StanceApplied  string          `json:"stance_applied,omitempty"`
}

// EvaluationConfig configures tactical stances and stochastic temperature sampling.
type EvaluationConfig struct {
	Stance      string     // balanced, aggressive, covert, diplomatic, inquisitive
	Temperature float64    // 0.0 = deterministic argmax; >0.0 = Boltzmann Softmax sampling
	Rnd         *rand.Rand // Optional PRNG source
}

// DefaultEvaluationConfig returns standard deterministic balanced evaluation settings.
func DefaultEvaluationConfig() EvaluationConfig {
	return EvaluationConfig{
		Stance:      "balanced",
		Temperature: 0.0,
	}
}

// EvaluateChoice selects the most fitting action using default balanced deterministic evaluation.
func (p *Persona) EvaluateChoice(sceneDescription string, candidates []CandidateAction) DecisionEvaluation {
	return p.EvaluateChoiceWithConfig(sceneDescription, candidates, DefaultEvaluationConfig())
}

// EvaluateChoiceWithConfig evaluates candidates applying tactical stance and temperature sampling.
func (p *Persona) EvaluateChoiceWithConfig(sceneDescription string, candidates []CandidateAction, cfg EvaluationConfig) DecisionEvaluation {
	if len(candidates) == 0 {
		return DecisionEvaluation{}
	}

	stance := strings.ToLower(strings.TrimSpace(cfg.Stance))
	if stance == "" {
		stance = "balanced"
	}

	scores := make([]float64, len(candidates))
	riskAssessments := make([]string, len(candidates))
	rationales := make([]string, len(candidates))

	goalWords := strings.ToLower(p.PsychologicalDrivers.PrimaryGoal)
	isAthletic := false
	for _, trait := range p.CanonicalTraits {
		tl := strings.ToLower(trait)
		if strings.Contains(tl, "athletic") || strings.Contains(tl, "decisive") {
			isAthletic = true
			break
		}
	}

	for i, cand := range candidates {
		score := 50.0 // Base neutral score

		// Match action description and keywords against character's primary goal
		for _, kw := range cand.Keywords {
			kwLower := strings.ToLower(kw)
			if strings.Contains(goalWords, kwLower) {
				score += 20.0
			}
			// Synergy with metaopinions
			for _, m := range p.Metaopinions {
				opinionLower := strings.ToLower(m.ExploratoryMetaopinion)
				if strings.Contains(opinionLower, kwLower) {
					score += (15.0 * m.CertaintyScore)
				}
			}
		}

		// Tactical Stance skewing
		switch stance {
		case "aggressive":
			if cand.TacticalType == "AGGRESSIVE" {
				score += 40.0
			} else {
				score -= 20.0
			}
			score -= cand.BaseRisk * 10.0 // Reduced fear of risk
		case "covert":
			if cand.TacticalType == "COVERT" {
				score += 40.0
			} else {
				score -= 20.0
			}
			if cand.BaseRisk < 0.4 {
				score += 15.0 // Rewarding stealth safety
			}
		case "diplomatic":
			if cand.TacticalType == "DIPLOMATIC" {
				score += 40.0
			} else {
				score -= 20.0
			}
		case "inquisitive":
			if cand.TacticalType == "INQUISITIVE" {
				score += 40.0
			} else {
				score -= 20.0
			}
		default: // balanced
			if cand.BaseRisk > 0.6 && !isAthletic {
				score -= cand.BaseRisk * 25.0
			} else {
				score += (1.0 - cand.BaseRisk) * 10.0
			}
		}

		// Vulnerability interaction
		vulnLower := strings.ToLower(p.PsychologicalDrivers.Vulnerability)
		for _, kw := range cand.Keywords {
			if strings.Contains(vulnLower, strings.ToLower(kw)) {
				score += 5.0
			}
		}

		riskDesc := "Manageable risk profile aligned with tactical doctrine."
		if cand.BaseRisk > 0.7 {
			riskDesc = "High operational hazard; exposes actor to immediate interception or betrayal."
		} else if cand.BaseRisk < 0.3 {
			riskDesc = "Low operational exposure; emphasizes covert preservation and strategic patience."
		}

		scores[i] = score
		riskAssessments[i] = riskDesc
		rationales[i] = p.formulateRationale(cand, score, stance)
	}

	// Deterministic selection if Temperature <= 0.001
	if cfg.Temperature <= 0.001 {
		bestIdx := 0
		bestScore := scores[0]
		for i := 1; i < len(scores); i++ {
			if scores[i] > bestScore {
				bestScore = scores[i]
				bestIdx = i
			}
		}
		return DecisionEvaluation{
			ChosenAction:   candidates[bestIdx],
			Score:          bestScore,
			Probability:    1.0,
			RiskAssessment: riskAssessments[bestIdx],
			Rationale:      rationales[bestIdx],
			StanceApplied:  stance,
		}
	}

	// Softmax / Boltzmann Temperature Sampling
	maxScore := scores[0]
	for _, s := range scores {
		if s > maxScore {
			maxScore = s
		}
	}

	expVals := make([]float64, len(scores))
	sumExp := 0.0
	for i, s := range scores {
		val := math.Exp((s - maxScore) / cfg.Temperature)
		expVals[i] = val
		sumExp += val
	}

	probs := make([]float64, len(scores))
	for i := range expVals {
		probs[i] = expVals[i] / sumExp
	}

	rnd := cfg.Rnd
	if rnd == nil {
		rnd = rand.New(rand.NewSource(1))
	}
	r := rnd.Float64()

	cum := 0.0
	chosenIdx := len(candidates) - 1
	for i, pVal := range probs {
		cum += pVal
		if r <= cum {
			chosenIdx = i
			break
		}
	}

	chosenRationale := fmt.Sprintf("%s [Stochastic sampling T=%.2f, P=%.2f%%]", rationales[chosenIdx], cfg.Temperature, probs[chosenIdx]*100.0)

	return DecisionEvaluation{
		ChosenAction:   candidates[chosenIdx],
		Score:          scores[chosenIdx],
		Probability:    probs[chosenIdx],
		RiskAssessment: riskAssessments[chosenIdx],
		Rationale:      chosenRationale,
		StanceApplied:  stance,
	}
}

func (p *Persona) formulateRationale(cand CandidateAction, score float64, stance string) string {
	var sb strings.Builder
	sb.WriteString(p.Name)
	sb.WriteString(" selects '")
	sb.WriteString(cand.Description)
	sb.WriteString("'")
	if stance != "balanced" {
		sb.WriteString(" (Stance skew: ")
		sb.WriteString(stance)
		sb.WriteString(")")
	}
	sb.WriteString(" based on primary directive: ")
	sb.WriteString(p.PsychologicalDrivers.PrimaryGoal)
	sb.WriteString(". ")
	if len(p.Metaopinions) > 0 {
		sb.WriteString("Reflecting metaopinion: ")
		sb.WriteString(p.Metaopinions[0].Topic)
		sb.WriteString(" (")
		sb.WriteString(p.Metaopinions[0].ExploratoryMetaopinion)
		sb.WriteString(").")
	}
	return sb.String()
}
