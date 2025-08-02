package storage

import (
	"awesomeProject/model"
	"gorm.io/gorm"
	"time"
)

type IDiaryStorage interface {
	CreateDiary(diary *model.Diary) error
	GetDiaryByID(id uint) (*model.Diary, error)
	GetDiaryByUserIDAndDate(userID uint, date time.Time) (*model.Diary, error)
	GetDiariesByUserID(userID uint) ([]*model.Diary, error)
	UpdateDiary(diary *model.Diary) error
	DeleteDiary(id uint) error
}

func NewDiaryStorage(db *gorm.DB) IDiaryStorage {
	return &diaryStorage{db: db}
}

type diaryStorage struct {
	db *gorm.DB
}

func (s *diaryStorage) CreateDiary(diary *model.Diary) error {
	return s.db.Create(diary).Error
}

func (s *diaryStorage) GetDiaryByID(id uint) (*model.Diary, error) {
	var diary model.Diary
	if err := s.db.First(&diary, id).Error; err != nil {
		return nil, err
	}
	return &diary, nil
}

func (s *diaryStorage) GetDiaryByUserIDAndDate(userID uint, date time.Time) (*model.Diary, error) {
	var diary model.Diary
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	if err := s.db.Where("user_id = ? AND date >= ? AND date < ?", userID, startOfDay, endOfDay).First(&diary).Error; err != nil {
		return nil, err
	}
	return &diary, nil
}

func (s *diaryStorage) GetDiariesByUserID(userID uint) ([]*model.Diary, error) {
	var diaries []*model.Diary
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&diaries).Error; err != nil {
		return nil, err
	}
	return diaries, nil
}

func (s *diaryStorage) UpdateDiary(diary *model.Diary) error {
	return s.db.Save(diary).Error
}

func (s *diaryStorage) DeleteDiary(id uint) error {
	return s.db.Delete(&model.Diary{}, id).Error
}
