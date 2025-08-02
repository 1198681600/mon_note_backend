package service

import (
	"awesomeProject/model"
	"awesomeProject/storage"
	"awesomeProject/utils"
	"errors"
	"time"
)

type IAuthService interface {
	SendVerificationCode(email string) error
	Login(req *LoginRequest) (*AuthResponse, error)
	ValidateToken(token string) (*model.User, error)
	Logout(token string) error
	UpdateProfile(userID uint, req *UpdateProfileRequest) (*model.User, error)
}

type authService struct {
	userStorage storage.IUserStorage
}

func NewAuthService(userStorage storage.IUserStorage) IAuthService {
	return &authService{
		userStorage: userStorage,
	}
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

type UpdateProfileRequest struct {
	Nickname   *string `json:"nickname"`
	Gender     *string `json:"gender"`
	Age        *int    `json:"age"`
	Profession *string `json:"profession"`
	Avatar     *string `json:"avatar"`
}

func (s *authService) SendVerificationCode(email string) error {
	verification := &model.EmailVerification{
		Email:            email,
		VerificationCode: "111111",
		IsUsed:           false,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	}

	return s.userStorage.CreateEmailVerification(verification)
}



func (s *authService) Login(req *LoginRequest) (*AuthResponse, error) {
	verification, err := s.userStorage.GetEmailVerification(req.Email, req.Code)
	if err != nil {
		return nil, errors.New("验证码无效或已过期")
	}

	user, err := s.userStorage.GetUserByEmail(req.Email)
	if err != nil {
		// 如果用户不存在，创建新用户
		user = &model.User{
			Email:           req.Email,
			IsEmailVerified: true,
		}
		
		err = s.userStorage.CreateUser(user)
		if err != nil {
			return nil, errors.New("创建用户失败")
		}
	}

	err = s.userStorage.MarkVerificationAsUsed(verification.ID)
	if err != nil {
		return nil, err
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

func (s *authService) ValidateToken(token string) (*model.User, error) {
	session, err := s.userStorage.GetUserSessionByToken(token)
	if err != nil {
		return nil, errors.New("token无效或已过期")
	}
	return &session.User, nil
}

func (s *authService) Logout(token string) error {
	return s.userStorage.DeleteUserSession(token)
}

func (s *authService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*model.User, error) {
	user, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Age != nil {
		user.Age = *req.Age
	}
	if req.Profession != nil {
		user.Profession = *req.Profession
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	err = s.userStorage.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}