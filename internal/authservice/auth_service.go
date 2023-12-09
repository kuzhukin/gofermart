package authservice

import (
	"errors"
	"fmt"
	"gophermart/internal/apiserver/handler"
	"gophermart/internal/authservice/cryptographer"
	"gophermart/internal/authservice/storage"
)

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
		if errors.Is(storage.ErrIsAlreadySaved, err) {
			return "", fmt.Errorf("user with login=%s already registred, err=%w", login, handler.ErrIsAlreadyRegistred)
		}

		return "", fmt.Errorf("save user's info login=%s, err=%w", login, err)
	}

	return key, nil
}

func (s *AuthService) Authorize(login string, password string) (string, error) {
	userInfo, err := s.authStorage.GetUserInfo(login)
	if err != nil {
		if errors.Is(err, storage.ErrIsNotContains) {
			return "", fmt.Errorf("wasn't registred, err=%w", handler.ErrIsNotAutorized)
		}

		return "", fmt.Errorf("get user info login=%s, err=%w", login, err)
	}

	if userInfo.Password != password {
		return "", fmt.Errorf("bad password, err=%w", handler.ErrIsNotAutorized)
	}

	return userInfo.Key, nil
}

func (s *AuthService) Check(userKey string) (string, error) {
	info, err := s.authStorage.GetUserInfoByKey(userKey)
	if errors.Is(err, storage.ErrIsNotContains) {
		return "", fmt.Errorf("wasn't registred, err=%w", handler.ErrIsNotAutorized)
	}

	return info.Login, nil
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
