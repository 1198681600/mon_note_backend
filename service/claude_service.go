package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type IClaudeService interface {
	AnalyzeDiaryEmotion(diaryContent, diaryDate string, userContext map[string]interface{}) (*DiaryEmotionResponse, error)
	AnalyzeWeeklyEmotion(weekStart string, diaryData []map[string]interface{}) (*WeeklyEmotionResponse, error)
}

type claudeService struct {
	apiURL string
	apiKey string
	client *http.Client
}

func NewClaudeService() IClaudeService {
	return &claudeService{
		apiURL: "https://api.gptsapi.net/v1/chat/completions",
		apiKey: "sk-8zU7c731868d93c819a720e252bf4c8620b5738449364Bcm",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type ClaudeRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type EmotionAnalysisRequest struct {
	UserData []EmotionData `json:"user_data"`
}

type EmotionData struct {
	Timestamp string `json:"timestamp"`
	Emotion   string `json:"emotion"`
	Intensity int    `json:"intensity"`
	Content   string `json:"content"`
}

type EmotionAnalysisResponse struct {
	DailyPattern    []HourlyEmotion `json:"daily_pattern"`
	Summary         string          `json:"summary"`
	Recommendations []string        `json:"recommendations"`
	TrendAnalysis   string          `json:"trend_analysis"`
}

type HourlyEmotion struct {
	Hour            int     `json:"hour"`
	DominantEmotion string  `json:"dominant_emotion"`
	AverageScore    float64 `json:"average_score"`
	Count           int     `json:"count"`
}

type DiaryEmotionResponse struct {
	Emotions           []EmotionDetail    `json:"emotions"`
	GradientSuggestion GradientSuggestion `json:"gradient_suggestion"`
	Summary            EmotionSummary     `json:"summary"`
	Insights           []string           `json:"insights"`
	Recommendations    []string           `json:"recommendations"`
}

type EmotionDetail struct {
	Emotion     string  `json:"emotion"`
	Intensity   float64 `json:"intensity"`
	TimePeriod  string  `json:"time_period"`
	Color       string  `json:"color"`
	Description string  `json:"description"`
}

type GradientSuggestion struct {
	Type      string `json:"type"`
	Reasoning string `json:"reasoning"`
}

type EmotionSummary struct {
	DominantEmotion    string  `json:"dominant_emotion"`
	EmotionalStability float64 `json:"emotional_stability"`
	MoodTrend          string  `json:"mood_trend"`
	EnergyLevel        string  `json:"energy_level"`
}

type WeeklyEmotionResponse struct {
	WeeklyPattern   map[string]DailyEmotion `json:"weekly_pattern"`
	Insights        []string                `json:"insights"`
	Recommendations []string                `json:"recommendations"`
	EmotionScore    float64                 `json:"emotion_score"`
	StabilityScore  float64                 `json:"stability_score"`
}

type DailyEmotion struct {
	Dominant  string  `json:"dominant"`
	Intensity float64 `json:"intensity"`
}

func (s *claudeService) AnalyzeDiaryEmotion(diaryContent, diaryDate string, userContext map[string]interface{}) (*DiaryEmotionResponse, error) {
	prompt := s.buildDiaryAnalysisPrompt(diaryContent, diaryDate, userContext)

	claudeReq := ClaudeRequest{
		Model: "claude-3-sonnet-20240229",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	reqBody, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(claudeResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from Claude API")
	}

	var result DiaryEmotionResponse
	if err := json.Unmarshal([]byte(claudeResp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Claude response: %v", err)
	}

	return &result, nil
}

func (s *claudeService) AnalyzeWeeklyEmotion(weekStart string, diaryData []map[string]interface{}) (*WeeklyEmotionResponse, error) {
	prompt := s.buildWeeklyAnalysisPrompt(weekStart, diaryData)

	claudeReq := ClaudeRequest{
		Model: "claude-3-sonnet-20240229",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	reqBody, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(claudeResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from Claude API")
	}

	var result WeeklyEmotionResponse
	if err := json.Unmarshal([]byte(claudeResp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Claude response: %v", err)
	}

	return &result, nil
}

func (s *claudeService) buildEmotionAnalysisPrompt(data []EmotionData) string {
	dataJSON, _ := json.Marshal(data)

	return fmt.Sprintf(`你是一个专业的情绪分析师。请分析以下用户的情绪数据，识别一天中情绪变化的模式。

用户数据：
%s

请按照以下JSON格式返回分析结果，不要包含任何其他文本：

{
  "daily_pattern": [
    {
      "hour": 9,
      "dominant_emotion": "积极",
      "average_score": 7.5,
      "count": 3
    }
  ],
  "summary": "用户的情绪整体趋势描述",
  "recommendations": ["建议1", "建议2"],
  "trend_analysis": "详细的趋势分析"
}

分析要求：
1. 按小时统计情绪分布和强度
2. 识别情绪高峰和低谷时段
3. 提供个性化的情绪管理建议
4. 分析情绪变化的可能原因`, string(dataJSON))
}

func (s *claudeService) buildDiaryAnalysisPrompt(diaryContent, diaryDate string, userContext map[string]interface{}) string {
	userContextStr := ""
	if userContext != nil {
		if age, ok := userContext["age"]; ok {
			userContextStr += fmt.Sprintf("年龄: %v\n", age)
		}
		if gender, ok := userContext["gender"]; ok {
			userContextStr += fmt.Sprintf("性别: %v\n", gender)
		}
		if profession, ok := userContext["profession"]; ok {
			userContextStr += fmt.Sprintf("职业: %v\n", profession)
		}
	}

	return fmt.Sprintf(`你是一个专业的情绪分析专家，专门分析用户的日记内容并提取情绪信息。

## 任务要求
分析以下日记内容，提取其中的情绪信息，并按指定JSON格式返回结果。

## 日记内容
%s

## 日记日期  
%s

## 用户背景（可选）
%s

## 情绪色彩映射表
请从以下预定义情绪中选择最匹配的：
- 开心: #FFD700 (金黄色)
- 快乐: #FF6B6B (珊瑚红)
- 兴奋: #FF8E53 (橙色)
- 平静: #4ECDC4 (青绿色)
- 放松: #45B7D1 (天蓝色)
- 满足: #96CEB4 (薄荷绿)
- 难过: #6C5CE7 (紫色)
- 沮丧: #74B9FF (蓝色)
- 焦虑: #FDCB6E (淡黄色)
- 紧张: #E17055 (橙红色)
- 愤怒: #D63031 (红色)
- 烦躁: #E84393 (粉红色)
- 感动: #A29BFE (淡紫色)
- 温暖: #FFB8B8 (粉色)
- 孤独: #636E72 (灰色)
- 迷茫: #B2BEC3 (浅灰色)

## 渐变类型建议规则
最多三种情绪

## 输出格式
请严格按照以下JSON格式返回，不要包含任何其他文字：

{
  "emotions": [
    {
      "emotion": "情绪名称",
      "intensity": 0.0-1.0之间的数值,
      "time_period": "上午/中午/下午/晚上",
      "color": "对应的颜色代码",
      "description": "该情绪的简短描述(20字以内)"
    }
  ],
  "gradient_suggestion": {
    "type": "radial/sweep/multiPoint/wave/diagonal",
    "reasoning": "选择此渐变类型的原因"
  },
  "summary": {
    "dominant_emotion": "主导情绪",
    "emotional_stability": 1-10的评分,
    "mood_trend": "整体趋势描述",
    "energy_level": "低/中等/高"
  },
  "insights": [
    "分析洞察1",
    "分析洞察2",
    "分析洞察3"
  ],
  "recommendations": [
    "建议1", 
    "建议2",
    "建议3"
  ]
}

## 注意事项
1. intensity值必须在0.0-1.0之间
2. emotion必须从预定义列表中选择
3. time_period只能是：上午、中午、下午、晚上
4. 每个数组最多包含5个元素
5. 描述要简洁准确，贴合中文表达习惯
6. 分析要专业但易懂，避免过于学术化的表述`, diaryContent, diaryDate, userContextStr)
}

func (s *claudeService) buildWeeklyAnalysisPrompt(weekStart string, diaryData []map[string]interface{}) string {
	diaryDataJSON, _ := json.Marshal(diaryData)

	return fmt.Sprintf(`你是一个专业的情绪分析专家，专门分析用户的一周情绪数据并提供趋势分析。

## 任务要求
分析以下一周的情绪数据，识别情绪模式和趋势，并按指定JSON格式返回结果。

## 一周开始日期
%s

## 一周情绪数据
%s

## 输出格式
请严格按照以下JSON格式返回，不要包含任何其他文字：

{
  "weekly_pattern": {
    "monday": {"dominant": "平静", "intensity": 0.6},
    "tuesday": {"dominant": "开心", "intensity": 0.8},
    "wednesday": {"dominant": "焦虑", "intensity": 0.7},
    "thursday": {"dominant": "满足", "intensity": 0.7},
    "friday": {"dominant": "兴奋", "intensity": 0.9},
    "saturday": {"dominant": "放松", "intensity": 0.8},
    "sunday": {"dominant": "温暖", "intensity": 0.8}
  },
  "insights": [
    "本周整体情绪波动属于正常范围",
    "周三出现焦虑情绪，可能与工作压力相关",
    "周末情绪明显改善，休息对你很重要"
  ],
  "recommendations": [
    "周三是你的情绪低谷期，建议安排轻松的活动",
    "保持现有的周末放松习惯",
    "可以在周中增加一些减压活动"
  ],
  "emotion_score": 7.4,
  "stability_score": 6.8
}

## 分析要求
1. 识别每天的主导情绪和强度
2. 分析一周内的情绪变化趋势
3. 找出情绪高峰和低谷的规律
4. 提供针对性的情绪管理建议
5. 给出整体情绪评分(1-10分)和稳定性评分(1-10分)`, weekStart, string(diaryDataJSON))
}
