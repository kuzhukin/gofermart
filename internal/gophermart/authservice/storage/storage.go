package storage

import (
	"errors"
	"sync"
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
	GetUserInfoByKey(key Key) (*UserInfo, error)
}

var _ Storage = &AuthStorage{}

var ErrIsAlreadySaved = errors.New("is already saved")
var ErrIsNotContains = errors.New("isn't contains")

type AuthStorage struct {
	// TODO: Change to db
	sync.Mutex
	login2Info map[Login]*UserInfo
	key2Info   map[Key]*UserInfo
}

func New() *AuthStorage {
	return &AuthStorage{
		login2Info: make(map[Login]*UserInfo),
		key2Info:   make(map[Key]*UserInfo),
	}
}

func (s *AuthStorage) SaveUserInfo(login Login, password Password, key Key) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.login2Info[login]
	if ok {
		return ErrIsAlreadySaved
	}

	userInfo := &UserInfo{Login: login, Password: password, Key: key}
	s.login2Info[login] = userInfo
	s.key2Info[key] = userInfo

	return nil
}

func (s *AuthStorage) GetUserInfo(login string) (*UserInfo, error) {
	s.Lock()
	defer s.Unlock()

	data, ok := s.login2Info[login]
	if !ok {
		return nil, ErrIsNotContains
	}

	return data, nil
}

func (s *AuthStorage) GetUserInfoByKey(key string) (*UserInfo, error) {
	s.Lock()
	defer s.Unlock()

	data, ok := s.key2Info[key]
	if !ok {
		return nil, ErrIsNotContains
	}

	return data, nil
}
