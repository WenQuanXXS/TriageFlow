package service

import (
	"testing"

	"github.com/triageflow/backend/model"
)

func TestRuleEngine_ChestPain(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "突发胸痛，持续半小时"}, nil)
	if result.RuleTriggered != "chest_pain" {
		t.Errorf("expected chest_pain rule, got '%s'", result.RuleTriggered)
	}
	if result.FinalPriority != "urgent" {
		t.Errorf("expected urgent, got '%s'", result.FinalPriority)
	}
	if result.FinalDepartment != "Emergency" {
		t.Errorf("expected Emergency, got '%s'", result.FinalDepartment)
	}
}

func TestRuleEngine_BreathingDifficulty(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "呼吸困难，无法平躺"}, nil)
	if result.RuleTriggered != "breathing_difficulty" {
		t.Errorf("expected breathing_difficulty, got '%s'", result.RuleTriggered)
	}
}

func TestRuleEngine_RiskSignalFromAI(t *testing.T) {
	re := NewRuleEngine()
	aiResult := &model.TriageResult{
		RiskSignals: []string{"疑似胸闷"},
	}
	result := re.Evaluate(PatientInfo{ChiefComplaint: "感觉不舒服"}, aiResult)
	if result.RuleTriggered != "chest_pain" {
		t.Errorf("expected chest_pain from risk signal, got '%s'", result.RuleTriggered)
	}
}

func TestRuleEngine_NoTrigger(t *testing.T) {
	re := NewRuleEngine()
	aiResult := &model.TriageResult{
		SuggestedPri:   "normal",
		CandidateDepts: []string{"Neurology"},
	}
	result := re.Evaluate(PatientInfo{ChiefComplaint: "头痛三天"}, aiResult)
	if result.RuleTriggered != "" {
		t.Errorf("expected no rule trigger, got '%s'", result.RuleTriggered)
	}
	if result.FinalPriority != "normal" {
		t.Errorf("expected normal, got '%s'", result.FinalPriority)
	}
	if result.FinalDepartment != "Neurology" {
		t.Errorf("expected Neurology, got '%s'", result.FinalDepartment)
	}
}

func TestRuleEngine_Seizure(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "孩子突然抽搐"}, nil)
	if result.RuleTriggered != "seizure" {
		t.Errorf("expected seizure, got '%s'", result.RuleTriggered)
	}
}

func TestRuleEngine_SevereAllergy(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "吃了海鲜后严重过敏"}, nil)
	if result.RuleTriggered != "severe_allergy" {
		t.Errorf("expected severe_allergy, got '%s'", result.RuleTriggered)
	}
}

func TestRuleEngine_Stroke(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "老人突然口齿不清，一侧肢体无力"}, nil)
	if result.RuleTriggered != "stroke" {
		t.Errorf("expected stroke, got '%s'", result.RuleTriggered)
	}
	if result.FinalPriority != "urgent" {
		t.Errorf("expected urgent, got '%s'", result.FinalPriority)
	}
	if result.FinalDepartment != "Emergency" {
		t.Errorf("expected Emergency, got '%s'", result.FinalDepartment)
	}
}

func TestRuleEngine_AcuteAbdomen(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "剧烈腹痛两小时，越来越严重"}, nil)
	if result.RuleTriggered != "acute_abdomen" {
		t.Errorf("expected acute_abdomen, got '%s'", result.RuleTriggered)
	}
	if result.FinalPriority != "high" {
		t.Errorf("expected high, got '%s'", result.FinalPriority)
	}
	if result.FinalDepartment != "General Surgery" {
		t.Errorf("expected General Surgery, got '%s'", result.FinalDepartment)
	}
}

func TestRuleEngine_Poisoning(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "孩子误食了清洁剂"}, nil)
	if result.RuleTriggered != "poisoning" {
		t.Errorf("expected poisoning, got '%s'", result.RuleTriggered)
	}
	if result.FinalDepartment != "Emergency" {
		t.Errorf("expected Emergency, got '%s'", result.FinalDepartment)
	}
}

func TestRuleEngine_SevereTrauma(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "车祸后多处受伤"}, nil)
	if result.RuleTriggered != "severe_trauma" {
		t.Errorf("expected severe_trauma, got '%s'", result.RuleTriggered)
	}
	if result.FinalPriority != "urgent" {
		t.Errorf("expected urgent, got '%s'", result.FinalPriority)
	}
}

func TestRuleEngine_HighFever(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "持续高热三天不退"}, nil)
	if result.RuleTriggered != "high_fever" {
		t.Errorf("expected high_fever, got '%s'", result.RuleTriggered)
	}
	if result.FinalPriority != "high" {
		t.Errorf("expected high, got '%s'", result.FinalPriority)
	}
	if result.FinalDepartment != "Internal Medicine" {
		t.Errorf("expected Internal Medicine, got '%s'", result.FinalDepartment)
	}
}

func TestRuleEngine_VaginalBleeding(t *testing.T) {
	re := NewRuleEngine()
	result := re.Evaluate(PatientInfo{ChiefComplaint: "非经期阴道出血"}, nil)
	if result.RuleTriggered != "vaginal_bleeding" {
		t.Errorf("expected vaginal_bleeding, got '%s'", result.RuleTriggered)
	}
	if result.FinalDepartment != "Gynecology" {
		t.Errorf("expected Gynecology, got '%s'", result.FinalDepartment)
	}
}
