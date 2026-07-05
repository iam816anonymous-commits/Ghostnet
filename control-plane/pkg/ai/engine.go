package ai

import (
	"fmt"
	"ghostnet/control-plane/pkg/common"
)

type Engine struct {
	Goal string // "speed", "balanced", "privacy"
}

func NewEngine(goal string) *Engine {
	return &Engine{Goal: goal}
}

func (e *Engine) RecommendStrategy(riskScore float64) string {
	if riskScore > 7.0 {
		return "Switch to Maximum Privacy mode and rotate identities."
	}
	if e.Goal == "speed" {
		return "Optimize for low-latency routes; minimize padding."
	}
	return "Maintain balanced operation with moderate cover traffic."
}

func (e *Engine) EstimateRisk(activeUsers int, mode common.PrivacyMode) float64 {
	risk := 10.0

	switch mode {
	case common.ModeMaximumPrivacy:
		risk -= 8.0
	case common.ModeHigh:
		risk -= 5.0
	case common.ModeMedium:
		risk -= 3.0
	case common.ModeLow:
		risk -= 1.0
	}

	if activeUsers > 1000 {
		risk -= 1.0
	}
	if activeUsers > 10000 {
		risk -= 0.5
	}

	if risk < 0 {
		risk = 0
	}
	return risk
}

func (e *Engine) ExplainDecision(metric string, value interface{}) string {
	return fmt.Sprintf("AI Decision: Adjusting %s to %v based on current network conditions and user goal (%s).", metric, value, e.Goal)
}
