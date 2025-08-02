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
	AnalyzeDailyEmotionPattern(data []EmotionData) (*EmotionAnalysisResponse, error)
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
	DailyPattern     []HourlyEmotion `json:"daily_pattern"`
	Summary          string          `json:"summary"`
	Recommendations  []string        `json:"recommendations"`
	TrendAnalysis    string          `json:"trend_analysis"`
}

type HourlyEmotion struct {
	Hour            int     `json:"hour"`
	DominantEmotion string  `json:"dominant_emotion"`
	AverageScore    float64 `json:"average_score"`
	Count           int     `json:"count"`
}

func (s *claudeService) AnalyzeDailyEmotionPattern(data []EmotionData) (*EmotionAnalysisResponse, error) {
	prompt := s.buildEmotionAnalysisPrompt(data)
	
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

	var result EmotionAnalysisResponse
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