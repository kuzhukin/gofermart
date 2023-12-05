package storage

import (
	"errors"
)

type Login = string
type Password = string
type Key = string

type UserInfo struct {
	Login    Login
	Password Password
	Key      Key
}

type Storage interface {
	SaveUserInfo(login Login, password Password, key Key) error
	GetUserInfo(login Login) (*UserInfo, error)
}

var _ Storage = &AuthStorage{}

var ErrIsAlreadySaved = errors.New("is already saved")
var ErrIsNotContains = errors.New("isn't contains")

type AuthStorage struct {
	// TODO: Change to db
	db map[string]*UserInfo
}

func New() *AuthStorage {
	return &AuthStorage{
		db: make(map[Login]*UserInfo),
	}
}

func (s *AuthStorage) SaveUserInfo(login Login, password Password, key Key) error {
	_, ok := s.db[login]
	if ok {
		return ErrIsAlreadySaved
	}

	s.db[login] = &UserInfo{Login: login, Password: password, Key: key}

	return nil
}

func (s *AuthStorage) GetUserInfo(login string) (*UserInfo, error) {
	data, ok := s.db[login]
	if !ok {
		return nil, ErrIsNotContains
	}

	return data, nil
}
