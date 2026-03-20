package service

import (
	"strings"

	"github.com/triageflow/backend/model"
)

type Rule struct {
	Name        string
	Keywords    []string
	ForcePri    string
	ForceDept   string
	Description string
}

type RuleEngine struct {
	rules []Rule
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules: []Rule{
			// Cardiac emergencies
			{Name: "chest_pain", Keywords: []string{"胸痛", "胸闷", "心绞痛", "chest pain"}, ForcePri: "urgent", ForceDept: "Emergency", Description: "胸痛可能提示心肌梗死，需紧急处理"},
			// Respiratory emergencies
			{Name: "breathing_difficulty", Keywords: []string{"呼吸困难", "喘不上气", "呼吸急促", "breathing difficulty"}, ForcePri: "urgent", ForceDept: "Emergency", Description: "呼吸困难需紧急评估"},
			// Neurological emergencies
			{Name: "consciousness_disorder", Keywords: []string{"意识障碍", "昏迷", "意识不清", "晕厥"}, ForcePri: "urgent", ForceDept: "Emergency", Description: "意识障碍需紧急处理"},
			{Name: "stroke", Keywords: []string{"偏瘫", "口齿不清", "一侧肢体无力", "面瘫", "中风"}, ForcePri: "urgent", ForceDept: "Emergency", Description: "疑似脑卒中，需紧急溶栓评估"},
			// Hemorrhage
			{Name: "major_bleeding", Keywords: []string{"大出血", "大量出血", "出血不止"}, ForcePri: "urgent", ForceDept: "Emergency", Description: "大出血需紧急止血处理"},
			// Seizure
			{Name: "seizure", Keywords: []string{"抽搐", "癫痫发作", "惊厥"}, ForcePri: "urgent", ForceDept: "Emergency", Description: "抽搐需紧急处理"},
			// Allergy
			{Name: "severe_allergy", Keywords: []string{"严重过敏", "过敏性休克", "喉头水肿"}, ForcePri: "urgent", ForceDept: "Emergency", Description: "严重过敏反应需紧急处理"},
			// Acute abdomen
			{Name: "acute_abdomen", Keywords: []string{"剧烈腹痛", "腹部剧痛", "阑尾炎"}, ForcePri: "high", ForceDept: "General Surgery", Description: "急腹症需外科紧急评估"},
			// Poisoning
			{Name: "poisoning", Keywords: []string{"中毒", "误食", "服毒"}, ForcePri: "urgent", ForceDept: "Emergency", Description: "中毒需紧急洗胃或解毒处理"},
			// Severe trauma
			{Name: "severe_trauma", Keywords: []string{"车祸", "高处坠落", "严重外伤"}, ForcePri: "urgent", ForceDept: "Emergency", Description: "严重创伤需紧急救治"},
			// High fever
			{Name: "high_fever", Keywords: []string{"高热", "高烧"}, ForcePri: "high", ForceDept: "Internal Medicine", Description: "高热需尽快退热并排查感染源"},
			// Vaginal bleeding (non-menstrual)
			{Name: "vaginal_bleeding", Keywords: []string{"阴道出血", "阴道大量出血"}, ForcePri: "high", ForceDept: "Gynecology", Description: "异常阴道出血需紧急妇科评估"},
		},
	}
}

// Evaluate checks the chief complaint and AI risk signals against safety rules.
// If a rule triggers, it overrides the AI-suggested priority and department.
func (re *RuleEngine) Evaluate(chiefComplaint string, aiResult *model.TriageResult) *model.RuleEngineResult {
	complaint := strings.ToLower(chiefComplaint)

	for _, rule := range re.rules {
		for _, kw := range rule.Keywords {
			if strings.Contains(complaint, strings.ToLower(kw)) {
				return &model.RuleEngineResult{
					RuleTriggered:   rule.Name,
					Reason:          rule.Description,
					FinalPriority:   rule.ForcePri,
					FinalDepartment: rule.ForceDept,
				}
			}
		}
	}

	// Also check AI-extracted risk signals
	if aiResult != nil {
		for _, signal := range aiResult.RiskSignals {
			signalLower := strings.ToLower(signal)
			for _, rule := range re.rules {
				for _, kw := range rule.Keywords {
					if strings.Contains(signalLower, strings.ToLower(kw)) {
						return &model.RuleEngineResult{
							RuleTriggered:   rule.Name,
							Reason:          rule.Description,
							FinalPriority:   rule.ForcePri,
							FinalDepartment: rule.ForceDept,
						}
					}
				}
			}
		}
	}

	// No rule triggered — pass through AI suggestion
	dept := ""
	pri := "normal"
	if aiResult != nil {
		pri = aiResult.SuggestedPri
		if len(aiResult.CandidateDepts) > 0 {
			dept = aiResult.CandidateDepts[0]
		}
	}
	return &model.RuleEngineResult{
		FinalPriority:   pri,
		FinalDepartment: dept,
	}
}
