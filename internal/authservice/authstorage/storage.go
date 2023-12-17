package authstorage

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/sql"
)

type User = sql.User

type Storage interface {
	SaveUser(ctx context.Context, login string, token string) error
	GetUser(ctx context.Context, login string) (*User, error)
	GetUserByToken(ctx context.Context, token string) (*User, error)
}

var _ Storage = &AuthStorageImpl{}

var ErrIsAlreadySaved = errors.New("is already saved")
var ErrIsNotContains = sql.ErrUserIsNotFound

type AuthStorageImpl struct {
	sqlController *sql.Controller
}

func New(ctrl *sql.Controller) *AuthStorageImpl {
	return &AuthStorageImpl{
		sqlController: ctrl,
	}
}

func (s *AuthStorageImpl) SaveUser(ctx context.Context, login string, token string) error {
	_, err := s.sqlController.FindUser(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrUserIsNotFound) {
			return s.sqlController.CreateUser(ctx, login, token)
		}

		return fmt.Errorf("find user=%s, err=%w", login, err)
	}

	return ErrIsAlreadySaved
}

func (s *AuthStorageImpl) GetUser(ctx context.Context, login string) (*User, error) {
	user, err := s.sqlController.FindUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("find user=%s, err=%w", login, err)
	}

	return user, nil
}

func (s *AuthStorageImpl) GetUserByToken(ctx context.Context, token string) (*User, error) {
	user, err := s.sqlController.FindUserByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("find user by token=%s, err=%w", token, err)
	}

	return user, nil
}
