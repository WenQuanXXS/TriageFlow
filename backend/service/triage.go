package service

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/triageflow/backend/model"
)

// PatientInfo carries all patient-submitted information for triage.
type PatientInfo struct {
	ChiefComplaint   string
	Age              int
	Gender           string
	Temperature      float64
	PainLevel        int
	SpecialCondition string
}

// Triager defines the interface for triage services.
// Implement this interface to swap between mock and real LLM.
type Triager interface {
	PerformTriage(ctx context.Context, info PatientInfo) (*model.TriageResult, string, error)
}

// symptom-to-department mapping for mock
var symptomDeptMap = map[string]string{
	// Neurology — 神经内科
	"头痛": "Neurology", "头晕": "Neurology", "头疼": "Neurology",
	"偏头痛": "Neurology", "面瘫": "Neurology", "手脚麻木": "Neurology",

	// Cardiology — 心内科
	"胸痛": "Cardiology", "胸闷": "Cardiology", "心悸": "Cardiology",
	"心慌": "Cardiology", "心绞痛": "Cardiology",

	// Pulmonology — 呼吸内科
	"咳嗽": "Pulmonology", "咳痰": "Pulmonology", "气喘": "Pulmonology",
	"哮喘": "Pulmonology", "感冒": "Pulmonology",
	"鼻塞": "Pulmonology", "流涕": "Pulmonology", "咽痛": "Pulmonology",

	// Gastroenterology — 消化内科
	"腹痛": "Gastroenterology", "肚子痛": "Gastroenterology", "胃痛": "Gastroenterology",
	"恶心": "Gastroenterology", "呕吐": "Gastroenterology",
	"腹泻": "Gastroenterology", "便秘": "Gastroenterology", "便血": "Gastroenterology",
	"胃胀": "Gastroenterology", "反酸": "Gastroenterology", "烧心": "Gastroenterology",

	// Orthopedics — 骨科
	"骨折": "Orthopedics", "摔伤": "Orthopedics",
	"关节痛": "Orthopedics", "关节肿": "Orthopedics",
	"腰痛": "Orthopedics", "腰酸": "Orthopedics",
	"颈椎痛": "Orthopedics", "颈椎": "Orthopedics", "扭伤": "Orthopedics",

	// General Surgery — 普外科
	"外伤": "General Surgery", "割伤": "General Surgery",
	"烧伤": "General Surgery", "烫伤": "General Surgery",
	"肿块": "General Surgery", "疝气": "General Surgery", "阑尾炎": "General Surgery",

	// Dermatology — 皮肤科
	"皮疹": "Dermatology", "瘙痒": "Dermatology", "湿疹": "Dermatology",
	"荨麻疹": "Dermatology", "痤疮": "Dermatology", "脱发": "Dermatology",

	// Pediatrics — 儿科
	"小儿": "Pediatrics", "孩子": "Pediatrics", "儿童": "Pediatrics",
	"婴儿": "Pediatrics", "宝宝": "Pediatrics",

	// Ophthalmology — 眼科
	"眼睛痛": "Ophthalmology", "眼痛": "Ophthalmology", "视力下降": "Ophthalmology",
	"视力模糊": "Ophthalmology", "眼红": "Ophthalmology", "眼睛干涩": "Ophthalmology",
	"流泪": "Ophthalmology", "眼睛": "Ophthalmology",

	// ENT — 耳鼻喉科
	"耳鸣": "ENT", "耳痛": "ENT", "听力下降": "ENT",
	"鼻炎": "ENT", "鼻出血": "ENT", "鼻子": "ENT",
	"喉咙痛": "ENT", "声音嘶哑": "ENT", "吞咽困难": "ENT",

	// Stomatology — 口腔科
	"牙痛": "Stomatology", "牙龈出血": "Stomatology", "牙龈肿": "Stomatology",
	"口腔溃疡": "Stomatology", "智齿": "Stomatology",

	// Urology — 泌尿外科
	"尿频": "Urology", "尿急": "Urology", "尿痛": "Urology",
	"血尿": "Urology", "肾结石": "Urology", "排尿困难": "Urology",

	// Gynecology — 妇科
	"月经不调": "Gynecology", "痛经": "Gynecology", "白带异常": "Gynecology",
	"阴道出血": "Gynecology",

	// Endocrinology — 内分泌科
	"糖尿病": "Endocrinology", "甲状腺": "Endocrinology",
	"多饮": "Endocrinology", "多尿": "Endocrinology", "甲亢": "Endocrinology",

	// Psychiatry — 精神心理科
	"失眠": "Psychiatry", "焦虑": "Psychiatry", "抑郁": "Psychiatry",
	"幻觉": "Psychiatry", "烦躁": "Psychiatry",

	// Rheumatology — 风湿免疫科
	"风湿": "Rheumatology", "类风湿": "Rheumatology",
	"红斑": "Rheumatology", "关节晨僵": "Rheumatology",

	// Internal Medicine — 内科 (general fallback symptoms)
	"发烧": "Internal Medicine", "发热": "Internal Medicine",
	"乏力": "Internal Medicine", "高热": "Internal Medicine",
	"过敏": "Internal Medicine",
}

// high-risk keywords for mock risk signal detection
var riskKeywords = []string{
	"胸痛", "胸闷", "呼吸困难", "喘不上气", "意识障碍", "昏迷", "晕厥",
	"大出血", "出血不止", "抽搐", "癫痫", "惊厥", "严重过敏", "过敏性休克",
	"高热", "剧烈疼痛",
	"中毒", "误食", "服毒",
	"偏瘫", "口齿不清", "一侧肢体无力",
	"剧烈腹痛", "腹部剧痛",
	"车祸", "高处坠落", "严重外伤",
}

// all symptom keywords for extraction
var symptomKeywords = []string{
	// 神经
	"头痛", "头晕", "头疼", "偏头痛", "面瘫", "手脚麻木",
	// 心血管
	"胸痛", "胸闷", "心悸", "心慌", "心绞痛",
	// 呼吸
	"咳嗽", "咳痰", "气喘", "哮喘", "感冒", "鼻塞", "流涕", "咽痛",
	"呼吸困难", "喘不上气",
	// 消化
	"腹痛", "肚子痛", "胃痛", "恶心", "呕吐",
	"腹泻", "便秘", "便血", "胃胀", "反酸", "烧心",
	// 骨科
	"骨折", "摔伤", "关节痛", "关节肿", "腰痛", "腰酸",
	"颈椎痛", "颈椎", "扭伤",
	// 外科
	"外伤", "割伤", "烧伤", "烫伤", "肿块", "疝气", "阑尾炎",
	// 皮肤
	"皮疹", "瘙痒", "湿疹", "荨麻疹", "痤疮", "脱发", "红肿",
	// 眼科
	"眼睛痛", "眼痛", "视力下降", "视力模糊", "眼红", "眼睛干涩", "流泪",
	// 耳鼻喉
	"耳鸣", "耳痛", "听力下降", "鼻炎", "鼻出血",
	"喉咙痛", "声音嘶哑", "吞咽困难",
	// 口腔
	"牙痛", "牙龈出血", "牙龈肿", "口腔溃疡", "智齿",
	// 泌尿
	"尿频", "尿急", "尿痛", "血尿", "肾结石", "排尿困难",
	// 妇科
	"月经不调", "痛经", "白带异常", "阴道出血",
	// 内分泌
	"糖尿病", "甲状腺", "多饮", "多尿", "甲亢",
	// 精神心理
	"失眠", "焦虑", "抑郁", "幻觉", "烦躁",
	// 风湿
	"风湿", "类风湿", "红斑", "关节晨僵",
	// 全身/急症
	"发烧", "发热", "高热", "乏力", "疼痛",
	"抽搐", "癫痫", "惊厥",
	"昏迷", "晕厥", "意识不清",
	"大出血", "出血不止",
	"过敏", "严重过敏",
	"中毒", "误食",
	"剧烈腹痛", "偏瘫", "口齿不清",
}

// MockTriageService simulates LLM triage using keyword matching.
type MockTriageService struct{}

func NewMockTriageService() *MockTriageService {
	return &MockTriageService{}
}

func (s *MockTriageService) PerformTriage(_ context.Context, info PatientInfo) (*model.TriageResult, string, error) {
	result := &model.TriageResult{}

	complaint := info.ChiefComplaint

	// Extract symptoms
	for _, kw := range symptomKeywords {
		if strings.Contains(complaint, kw) {
			result.Symptoms = append(result.Symptoms, kw)
		}
	}
	if len(result.Symptoms) == 0 {
		result.Symptoms = []string{"未识别到明确症状"}
	}

	// Detect risk signals
	for _, kw := range riskKeywords {
		if strings.Contains(complaint, kw) {
			result.RiskSignals = append(result.RiskSignals, kw)
		}
	}

	// Map to departments
	deptSet := map[string]bool{}
	for _, kw := range result.Symptoms {
		if dept, ok := symptomDeptMap[kw]; ok {
			deptSet[dept] = true
		}
	}
	for dept := range deptSet {
		result.CandidateDepts = append(result.CandidateDepts, dept)
	}
	sort.Strings(result.CandidateDepts)
	if len(result.CandidateDepts) == 0 {
		result.CandidateDepts = []string{"Internal Medicine"}
	}

	// Determine priority
	if len(result.RiskSignals) > 0 {
		result.SuggestedPri = "high"
	} else if len(result.Symptoms) > 2 {
		result.SuggestedPri = "high"
	} else {
		result.SuggestedPri = "normal"
	}

	result.Reasoning = "基于主诉关键词匹配的模拟分诊结果（Mock）"

	// Generate raw output for audit
	raw, _ := json.Marshal(result)

	return result, string(raw), nil
}
