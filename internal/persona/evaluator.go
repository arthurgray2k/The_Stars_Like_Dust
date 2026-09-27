package persona

import (
	"math"
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
	RiskAssessment string          `json:"risk_assessment"`
	Rationale      string          `json:"rationale"`
}

// EvaluateChoice selects the most fitting action for a persona given the scene context.
func (p *Persona) EvaluateChoice(sceneDescription string, candidates []CandidateAction) DecisionEvaluation {
	if len(candidates) == 0 {
		return DecisionEvaluation{}
	}

	bestScore := -math.MaxFloat64
	var bestAction CandidateAction
	var bestRisk string
	var bestRationale string

	goalWords := strings.ToLower(p.PsychologicalDrivers.PrimaryGoal)

	for _, cand := range candidates {
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

		// Penalize high risk unless character is aggressive/athletic
		isAthletic := false
		for _, trait := range p.CanonicalTraits {
			if strings.Contains(strings.ToLower(trait), "athletic") || strings.Contains(strings.ToLower(trait), "decisive") {
				isAthletic = true
				break
			}
		}

		if cand.BaseRisk > 0.6 && !isAthletic {
			score -= cand.BaseRisk * 25.0
		} else {
			score += (1.0 - cand.BaseRisk) * 10.0
		}

		// Vulnerability penalty
		vulnLower := strings.ToLower(p.PsychologicalDrivers.Vulnerability)
		for _, kw := range cand.Keywords {
			if strings.Contains(vulnLower, strings.ToLower(kw)) {
				// Vulnerability could either tempt or derail
				score += 5.0
			}
		}

		riskDesc := "Manageable risk profile aligned with tactical doctrine."
		if cand.BaseRisk > 0.7 {
			riskDesc = "High operational hazard; exposes actor to immediate interception or betrayal."
		} else if cand.BaseRisk < 0.3 {
			riskDesc = "Low operational exposure; emphasizes covert preservation and strategic patience."
		}

		rationale := p.formulateRationale(cand, score)

		if score > bestScore {
			bestScore = score
			bestAction = cand
			bestRisk = riskDesc
			bestRationale = rationale
		}
	}

	return DecisionEvaluation{
		ChosenAction:   bestAction,
		Score:          bestScore,
		RiskAssessment: bestRisk,
		Rationale:      bestRationale,
	}
}

func (p *Persona) formulateRationale(cand CandidateAction, score float64) string {
	var sb strings.Builder
	sb.WriteString(p.Name)
	sb.WriteString(" selects '")
	sb.WriteString(cand.Description)
	sb.WriteString("' based on core motivation: ")
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
