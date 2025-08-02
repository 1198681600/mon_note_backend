package model

import (
	"gorm.io/gorm"
	"time"
)

type Diary struct {
	gorm.Model
	UserID  uint      `gorm:"uniqueIndex:idx_user_date"`
	Date    time.Time `gorm:"uniqueIndex:idx_user_date"`
	Content string    `gorm:"type:text"`
}

func (Diary) TableName() string {
	return "diaries"
}
