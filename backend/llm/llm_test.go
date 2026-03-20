package llm

import (
	"context"
	"fmt"
	"testing"

	"github.com/triageflow/backend/config"
	"github.com/triageflow/backend/service"
)

func TestEinoTriage_Integration(t *testing.T) {
	cfg, err := config.Load("../config.json")
	if err != nil {
		t.Skipf("skipping: cannot load config: %v", err)
	}
	if !cfg.LLM.Enabled {
		t.Skip("skipping: LLM not enabled in config")
	}

	svc, err := NewEinoTriageService(&cfg.LLM)
	if err != nil {
		t.Fatalf("NewEinoTriageService: %v", err)
	}

	result, raw, err := svc.PerformTriage(context.Background(), service.PatientInfo{
		ChiefComplaint: "头痛三天，伴有恶心呕吐",
	})
	if err != nil {
		t.Fatalf("PerformTriage failed: %v", err)
	}

	fmt.Printf("Raw LLM output:\n%s\n\n", raw)
	fmt.Printf("Parsed result:\n")
	fmt.Printf("  Symptoms:       %v\n", result.Symptoms)
	fmt.Printf("  RiskSignals:    %v\n", result.RiskSignals)
	fmt.Printf("  CandidateDepts: %v\n", result.CandidateDepts)
	fmt.Printf("  SuggestedPri:   %s\n", result.SuggestedPri)
	fmt.Printf("  Reasoning:      %s\n", result.Reasoning)

	if len(result.Symptoms) == 0 {
		t.Error("expected at least one symptom")
	}
	if len(result.CandidateDepts) == 0 {
		t.Error("expected at least one candidate department")
	}
	if result.SuggestedPri == "" {
		t.Error("expected a suggested priority")
	}
}
