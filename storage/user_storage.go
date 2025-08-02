package storage

import (
	"awesomeProject/model"
	"time"
	"gorm.io/gorm"
)

type UserStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *UserStorage) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) UpdateUser(user *model.User) error {
	return s.db.Save(user).Error
}

func (s *UserStorage) CreateEmailVerification(verification *model.EmailVerification) error {
	return s.db.Create(verification).Error
}

func (s *UserStorage) GetEmailVerification(email, code string) (*model.EmailVerification, error) {
	var verification model.EmailVerification
	err := s.db.Where("email = ? AND verification_code = ? AND is_used = ? AND expires_at > ?", 
		email, code, false, time.Now()).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (s *UserStorage) MarkVerificationAsUsed(id uint) error {
	return s.db.Model(&model.EmailVerification{}).Where("id = ?", id).Update("is_used", true).Error
}

func (s *UserStorage) CreateUserSession(session *model.UserSession) error {
	return s.db.Create(session).Error
}

func (s *UserStorage) GetUserSessionByToken(token string) (*model.UserSession, error) {
	var session model.UserSession
	err := s.db.Preload("User").Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *UserStorage) DeleteUserSession(token string) error {
	return s.db.Where("token = ?", token).Delete(&model.UserSession{}).Error
}