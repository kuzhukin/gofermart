package authservice

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/apiserver/handler"
	"gophermart/internal/authservice/cryptographer"
	"gophermart/internal/storage/userstorage"
)

type AuthService struct {
	userStorage   userstorage.Storage
	cryptographer cryptographer.Cryptographer
}

func NewAuthService(storage userstorage.Storage, cryptographer cryptographer.Cryptographer) *AuthService {
	return &AuthService{
		userStorage:   storage,
		cryptographer: cryptographer,
	}
}

func (s *AuthService) Register(ctx context.Context, login string, password string) (string, error) {
	key, err := s.calcUserKey(login, password)
	if err != nil {
		return "", err
	}

	if err := s.userStorage.SaveUser(ctx, login, key); err != nil {
		if errors.Is(userstorage.ErrIsAlreadySaved, err) {
			return "", fmt.Errorf("user with login=%s already registred, err=%w", login, handler.ErrIsAlreadyRegistred)
		}

		return "", fmt.Errorf("save user's info login=%s, err=%w", login, err)
	}

	return key, nil
}

func (s *AuthService) Authorize(ctx context.Context, login string, password string) (string, error) {
	userInfo, err := s.userStorage.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, userstorage.ErrIsNotContains) {
			return "", fmt.Errorf("wasn't registred, err=%w", handler.ErrIsNotAutorized)
		}

		return "", fmt.Errorf("get user info login=%s, err=%w", login, err)
	}

	key, err := s.calcUserKey(login, password)
	if err != nil {
		return "", err
	}

	if userInfo.AuthToken != key {
		return "", fmt.Errorf("bad password, err=%w", handler.ErrIsNotAutorized)
	}

	return userInfo.AuthToken, nil
}

func (s *AuthService) Check(ctx context.Context, userKey string) (string, error) {
	info, err := s.userStorage.GetUserByToken(ctx, userKey)
	if err != nil {
		if errors.Is(err, userstorage.ErrIsNotContains) {
			return "", fmt.Errorf("wasn't registred, err=%w", handler.ErrIsNotAutorized)
		}

		return "", err
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
