package model

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Email           string         `gorm:"uniqueIndex;not null" json:"email"`
	IsEmailVerified bool           `gorm:"default:false" json:"is_email_verified"`
	Avatar          string         `gorm:"size:255" json:"avatar"`
	Nickname        string         `gorm:"size:50" json:"nickname"`
	Gender          string         `gorm:"size:10" json:"gender"`
	Age             int            `gorm:"default:0" json:"age"`
	Profession      string         `gorm:"size:100" json:"profession"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type EmailVerification struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Email            string    `gorm:"not null" json:"email"`
	VerificationCode string    `gorm:"not null" json:"verification_code"`
	IsUsed           bool      `gorm:"default:false" json:"is_used"`
	ExpiresAt        time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt        time.Time `json:"created_at"`
}

type UserSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}