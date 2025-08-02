package service

import (
	"awesomeProject/model"
	"awesomeProject/storage"
	"awesomeProject/utils"
	"errors"
	"time"
)

type AuthService struct {
	userStorage *storage.UserStorage
}

func NewAuthService(userStorage *storage.UserStorage) *AuthService {
	return &AuthService{
		userStorage: userStorage,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (s *AuthService) SendVerificationCode(email string) error {
	existingUser, _ := s.userStorage.GetUserByEmail(email)
	if existingUser != nil {
		return errors.New("用户已存在")
	}

	verification := &model.EmailVerification{
		Email:            email,
		VerificationCode: "111111",
		IsUsed:           false,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	}

	return s.userStorage.CreateEmailVerification(verification)
}

func (s *AuthService) Register(req *RegisterRequest) error {
	existingUser, _ := s.userStorage.GetUserByEmail(req.Email)
	if existingUser != nil {
		return errors.New("用户已存在")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:           req.Email,
		Password:        hashedPassword,
		IsEmailVerified: false,
	}

	return s.userStorage.CreateUser(user)
}

func (s *AuthService) VerifyEmail(req *VerifyEmailRequest) error {
	verification, err := s.userStorage.GetEmailVerification(req.Email, req.Code)
	if err != nil {
		return errors.New("验证码无效或已过期")
	}

	err = s.userStorage.MarkVerificationAsUsed(verification.ID)
	if err != nil {
		return err
	}

	user, err := s.userStorage.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.IsEmailVerified = true
	return s.userStorage.UpdateUser(user)
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	user, err := s.userStorage.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if !user.IsEmailVerified {
		return nil, errors.New("请先验证邮箱")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("密码错误")
	}

	token, err := utils.GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}

	session := &model.UserSession{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.userStorage.CreateUserSession(session)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) ValidateToken(token string) (*model.User, error) {
	session, err := s.userStorage.GetUserSessionByToken(token)
	if err != nil {
		return nil, errors.New("token无效或已过期")
	}
	return &session.User, nil
}

func (s *AuthService) Logout(token string) error {
	return s.userStorage.DeleteUserSession(token)
}