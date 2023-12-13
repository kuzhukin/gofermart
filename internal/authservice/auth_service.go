package authservice

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/apiserver/handler"
	"gophermart/internal/authservice/authstorage"
	"gophermart/internal/authservice/cryptographer"
)

type AuthService struct {
	authStorage   authstorage.Storage
	cryptographer cryptographer.Cryptographer
}

func NewAuthService(storage authstorage.Storage, cryptographer cryptographer.Cryptographer) *AuthService {
	return &AuthService{
		authStorage:   storage,
		cryptographer: cryptographer,
	}
}

func (s *AuthService) Register(ctx context.Context, login string, password string) (string, error) {
	key, err := s.calcUserKey(login, password)
	if err != nil {
		return "", err
	}

	if err := s.authStorage.SaveUser(ctx, login, key); err != nil {
		if errors.Is(authstorage.ErrIsAlreadySaved, err) {
			return "", fmt.Errorf("user with login=%s already registred, err=%w", login, handler.ErrIsAlreadyRegistred)
		}

		return "", fmt.Errorf("save user's info login=%s, err=%w", login, err)
	}

	return key, nil
}

func (s *AuthService) Authorize(ctx context.Context, login string, password string) (string, error) {
	userInfo, err := s.authStorage.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, authstorage.ErrIsNotContains) {
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
	info, err := s.authStorage.GetUserByToken(ctx, userKey)
	if errors.Is(err, authstorage.ErrIsNotContains) {
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
