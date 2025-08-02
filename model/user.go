package model

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Email           string         `gorm:"uniqueIndex;not null" json:"email"`
	Password        string         `gorm:"not null" json:"-"`
	IsEmailVerified bool           `gorm:"default:false" json:"is_email_verified"`
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