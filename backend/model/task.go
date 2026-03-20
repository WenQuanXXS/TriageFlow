package model

import "time"

type Task struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	PatientName      string    `json:"patient_name" gorm:"not null"`
	ChiefComplaint   string    `json:"chief_complaint" gorm:"not null"`
	Age              int       `json:"age"`
	Gender           string    `json:"gender"`
	Temperature      float64   `json:"temperature"`
	PainLevel        int       `json:"pain_level"`
	SpecialCondition string    `json:"special_condition" gorm:"type:text"`
	Status           string    `json:"status" gorm:"default:pending"`
	Priority         string    `json:"priority" gorm:"default:normal"`
	Department       string    `json:"department"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// LLM triage results
	Symptoms       string `json:"symptoms" gorm:"type:text"`
	RiskSignals    string `json:"risk_signals" gorm:"type:text"`
	CandidateDepts string `json:"candidate_depts" gorm:"type:text"`
	AISuggestedPri string `json:"ai_suggested_priority"`
	LLMRawOutput   string `json:"llm_raw_output" gorm:"type:text"`

	// Rule engine results
	RuleTriggered   string `json:"rule_triggered"`
	RuleReason      string `json:"rule_reason" gorm:"type:text"`
	FinalPriority   string `json:"final_priority"`
	FinalDepartment string `json:"final_department"`

	// Processing metadata
	TriageStatus string `json:"triage_status" gorm:"default:pending"`
}
