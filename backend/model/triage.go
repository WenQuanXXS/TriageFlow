package model

// TriageResult represents the structured output from the triage service (LLM or mock).
type TriageResult struct {
	Symptoms       []string `json:"symptoms"`
	RiskSignals    []string `json:"risk_signals"`
	CandidateDepts []string `json:"candidate_depts"`
	SuggestedPri   string   `json:"suggested_priority"`
	Reasoning      string   `json:"reasoning"`
}

// RuleEngineResult represents the rule engine's decision.
type RuleEngineResult struct {
	RuleTriggered   string `json:"rule_triggered"`
	Reason          string `json:"reason"`
	FinalPriority   string `json:"final_priority"`
	FinalDepartment string `json:"final_department"`
}
