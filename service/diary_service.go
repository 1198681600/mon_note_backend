package service

import (
	"awesomeProject/model"
	"awesomeProject/storage"
	"encoding/json"
	"fmt"
	"time"
)

type IDiaryService interface {
	CreateDiary(userID uint, content string, date time.Time) (*model.Diary, error)
	GetDiary(id, userID uint) (*model.Diary, error)
	GetDiaries(userID uint) ([]*model.Diary, error)
	UpdateDiary(id, userID uint, content string) (*model.Diary, error)
	DeleteDiary(id, userID uint) error
}

func NewDiaryService(storage storage.IDiaryStorage, claudeService IClaudeService) IDiaryService {
	return &diaryService{
		storage:      storage,
		claudeService: claudeService,
	}
}

type diaryService struct {
	storage      storage.IDiaryStorage
	claudeService IClaudeService
}

func (s *diaryService) CreateDiary(userID uint, content string, date time.Time) (*model.Diary, error) {
	// Check if a diary already exists for this user and date
	_, err := s.storage.GetDiaryByUserIDAndDate(userID, date)
	if err == nil {
		return nil, fmt.Errorf("diary already exists for this date")
	}

	// Generate emotion analysis
	emotionResult, err := s.claudeService.AnalyzeDiaryEmotion(content, date.Format("2006-01-02"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze emotion: %v", err)
	}

	// Convert emotion result to JSON
	emotionJSON, err := json.Marshal(emotionResult)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal emotion data: %v", err)
	}

	diary := &model.Diary{
		UserID:         userID,
		Content:        content,
		Date:           date,
		EmotionAnalysis: string(emotionJSON),
	}
	return diary, s.storage.CreateDiary(diary)
}

func (s *diaryService) GetDiary(id, userID uint) (*model.Diary, error) {
	diary, err := s.storage.GetDiaryByID(id)
	if err != nil {
		return nil, err
	}
	if diary.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	return diary, nil
}

func (s *diaryService) GetDiaries(userID uint) ([]*model.Diary, error) {
	return s.storage.GetDiariesByUserID(userID)
}

func (s *diaryService) UpdateDiary(id, userID uint, content string) (*model.Diary, error) {
	diary, err := s.GetDiary(id, userID)
	if err != nil {
		return nil, err
	}

	// Generate emotion analysis for updated content
	emotionResult, err := s.claudeService.AnalyzeDiaryEmotion(content, diary.Date.Format("2006-01-02"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze emotion: %v", err)
	}

	// Convert emotion result to JSON
	emotionJSON, err := json.Marshal(emotionResult)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal emotion data: %v", err)
	}

	diary.Content = content
	diary.EmotionAnalysis = string(emotionJSON)
	return diary, s.storage.UpdateDiary(diary)
}

func (s *diaryService) DeleteDiary(id, userID uint) error {
	if _, err := s.GetDiary(id, userID); err != nil {
		return err
	}
	return s.storage.DeleteDiary(id)
}
