package storage

import (
	"awesomeProject/model"
	"gorm.io/gorm"
	"time"
)

type IUserStorage interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateUser(user *model.User) error
	CreateEmailVerification(verification *model.EmailVerification) error
	GetEmailVerification(email, code string) (*model.EmailVerification, error)
	MarkVerificationAsUsed(id uint) error
	CreateUserSession(session *model.UserSession) error
	GetUserSessionByToken(token string) (*model.UserSession, error)
	DeleteUserSession(token string) error
}

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) IUserStorage {
	return &userStorage{db: db}
}

func (s *userStorage) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *userStorage) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userStorage) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userStorage) UpdateUser(user *model.User) error {
	return s.db.Save(user).Error
}

func (s *userStorage) CreateEmailVerification(verification *model.EmailVerification) error {
	return s.db.Create(verification).Error
}

func (s *userStorage) GetEmailVerification(email, code string) (*model.EmailVerification, error) {
	var verification model.EmailVerification
	err := s.db.Where("email = ? AND verification_code = ? AND is_used = ? AND expires_at > ?",
		email, code, false, time.Now()).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (s *userStorage) MarkVerificationAsUsed(id uint) error {
	return s.db.Model(&model.EmailVerification{}).Where("id = ?", id).Update("is_used", true).Error
}

func (s *userStorage) CreateUserSession(session *model.UserSession) error {
	return s.db.Create(session).Error
}

func (s *userStorage) GetUserSessionByToken(token string) (*model.UserSession, error) {
	var session model.UserSession
	err := s.db.Preload("User").Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *userStorage) DeleteUserSession(token string) error {
	return s.db.Where("token = ?", token).Delete(&model.UserSession{}).Error
}
