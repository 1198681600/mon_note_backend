package service

import (
	"awesomeProject/model"
	"awesomeProject/storage"
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

func NewDiaryService(storage storage.IDiaryStorage) IDiaryService {
	return &diaryService{storage: storage}
}

type diaryService struct {
	storage storage.IDiaryStorage
}

func (s *diaryService) CreateDiary(userID uint, content string, date time.Time) (*model.Diary, error) {
	// Check if a diary already exists for this user and date
	_, err := s.storage.GetDiaryByUserIDAndDate(userID, date)
	if err == nil {
		return nil, fmt.Errorf("diary already exists for this date")
	}

	diary := &model.Diary{
		UserID:  userID,
		Content: content,
		Date:    date,
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
	diary.Content = content
	return diary, s.storage.UpdateDiary(diary)
}

func (s *diaryService) DeleteDiary(id, userID uint) error {
	if _, err := s.GetDiary(id, userID); err != nil {
		return err
	}
	return s.storage.DeleteDiary(id)
}
