package storage

import (
	"errors"
	"gophermart/internal/gophermart/authservice/userinfo"
)

type Storage interface {
	SaveUserInfo(login string, password string, key userinfo.UserKey) error
	GetUserInfo(key userinfo.UserKey) (*userinfo.UserInfo, error)
}

var _ Storage = &AuthStorage{}

var (
	ErrAlreadyRegistred = errors.New("user is already registried")
)

type AuthStorage struct {
	// TODO: Change to db
	db map[userinfo.UserKey]userinfo.UserInfo
}

func New() *AuthStorage {
	return &AuthStorage{
		db: make(map[userinfo.UserKey]userinfo.UserInfo),
	}
}

func (s *AuthStorage) SaveUserInfo(login string, password string, key userinfo.UserKey) error {
	return nil
}

func (s *AuthStorage) GetUserInfo(key userinfo.UserKey) (*userinfo.UserInfo, error) {
	return nil, nil
}
