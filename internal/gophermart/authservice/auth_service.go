package authservice

import (
	"fmt"
	"gophermart/internal/gophermart/authservice/cryptographer"
	"gophermart/internal/gophermart/authservice/storage"
)

type UserAuthorizer interface {
	Authorize(login string, password string) (string, error)
}

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

func (s *AuthService) Register(login string, password string) (string, error) {
	key, err := s.calcUserKey(login, password)
	if err != nil {
		return "", err
	}

	if err := s.authStorage.SaveUserInfo(login, password, key); err != nil {
		return "", fmt.Errorf("save user's info login=%s, err=%w", login, err)
	}

	return key, nil
}

func (s *AuthService) calcUserKey(login string, password string) (string, error) {
	key, err := s.cryptographer.Encrypt(userDataToString(login, password))
	if err != nil {
		return "", fmt.Errorf("calc user key, err=%w", err)
	}

	return key, nil
}

func userDataToString(login, password string) string {
	return fmt.Sprintf("%s-%s", login, password)
}
