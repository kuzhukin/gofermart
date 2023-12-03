package authservice

import (
	"gophermart/internal/gophermart/authservice/cryptographer"
	"gophermart/internal/gophermart/authservice/storage"
	"gophermart/internal/gophermart/authservice/userinfo"
)

type UserRegistrator interface {
	Register(login string, password string) (userinfo.UserKey, error)
}

type UserAuthorizer interface {
	Authorize(login string, password string) (userinfo.UserKey, error)
}

var _ UserRegistrator = &AuthService{}

// var _ UserAuthorizer = &AuthService{}

type AuthService struct {
	authStorage   storage.Storage
	cryptographer cryptographer.Cryptographer
}

func NewAuthService(storage storage.Storage, cryptographer cryptographer.Cryptographer) *AuthService {
	return &AuthService{
		authStorage:   storage,
		cryptographer: cryptographer,
	}
}

func (s *AuthService) Register(login string, password string) (userinfo.UserKey, error) {
	user := userinfo.New(login, password)

	key, err := s.calcUserKey(user)
	if err != nil {
		return "", err
	}

	if err := s.authStorage.SaveUserInfo(user.Login, user.Password, key); err != nil {
		return "", err
	}

	return key, nil
}

func (s *AuthService) calcUserKey(user *userinfo.UserInfo) (userinfo.UserKey, error) {
	key, err := s.cryptographer.Encrypt(user.String())
	if err != nil {
		return "", nil
	}

	return userinfo.UserKey(key), nil
}
