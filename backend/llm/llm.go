package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	claude "github.com/cloudwego/eino-ext/components/model/claude"
	"github.com/cloudwego/eino/schema"
	"github.com/triageflow/backend/config"
	"github.com/triageflow/backend/model"
	"github.com/triageflow/backend/service"
)

// EinoTriageService implements Triager using EINO + Anthropic Claude API.
type EinoTriageService struct {
	cfg *config.LLMConfig
}

// NewEinoTriageService creates a new LLM triage service from config.
func NewEinoTriageService(cfg *config.LLMConfig) (*EinoTriageService, error) {
	if cfg.APIKey == "" || cfg.BaseURL == "" || cfg.Model == "" {
		return nil, fmt.Errorf("llm config requires api_key, base_url, and model")
	}
	return &EinoTriageService{cfg: cfg}, nil
}

func (s *EinoTriageService) PerformTriage(ctx context.Context, info service.PatientInfo) (*model.TriageResult, string, error) {
	maxTokens := 1024
	baseURL := s.cfg.BaseURL
	chatModel, err := claude.NewChatModel(ctx, &claude.Config{
		APIKey:    s.cfg.APIKey,
		BaseURL:   &baseURL,
		Model:     s.cfg.Model,
		MaxTokens: maxTokens,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create chat model: %w", err)
	}

	// Build user message with all patient info
	var parts []string
	parts = append(parts, fmt.Sprintf("患者主诉：%s", info.ChiefComplaint))
	if info.Age > 0 {
		parts = append(parts, fmt.Sprintf("年龄：%d岁", info.Age))
	}
	if info.Gender != "" {
		parts = append(parts, fmt.Sprintf("性别：%s", info.Gender))
	}
	if info.Temperature > 0 {
		parts = append(parts, fmt.Sprintf("体温：%.1f°C", info.Temperature))
	}
	if info.PainLevel > 0 {
		parts = append(parts, fmt.Sprintf("疼痛等级：%d/10", info.PainLevel))
	}
	if info.SpecialCondition != "" {
		parts = append(parts, fmt.Sprintf("特殊情况：%s", info.SpecialCondition))
	}
	userMessage := strings.Join(parts, "\n")

	messages := []*schema.Message{
		{Role: schema.System, Content: SystemPrompt},
		{Role: schema.User, Content: userMessage},
	}

	resp, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, "", fmt.Errorf("LLM generate failed: %w", err)
	}

	rawOutput := resp.Content

	// Strip possible markdown code fences
	content := strings.TrimSpace(rawOutput)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var result model.TriageResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, rawOutput, fmt.Errorf("failed to parse LLM response: %w (raw: %s)", err, rawOutput)
	}

	result.SuggestedPri = normalizePriority(result.SuggestedPri)

	return &result, rawOutput, nil
}

// normalizePriority maps LLM-returned priority to one of the valid values: urgent, high, normal.
func normalizePriority(pri string) string {
	switch strings.ToLower(strings.TrimSpace(pri)) {
	case "urgent", "emergency", "critical":
		return "urgent"
	case "high", "moderate", "medium":
		return "high"
	case "normal", "low", "routine":
		return "normal"
	default:
		return "normal"
	}
}
