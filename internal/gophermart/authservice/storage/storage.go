package storage

import (
	"errors"
)

type Login = string
type Password = string
type Key = string

type UserData struct {
	login     Login
	passsword Password
	key       Key
}

type Storage interface {
	SaveUserInfo(login Login, password Password, key Key) error
	GetUserInfo(login Login) (*UserData, error)
}

var _ Storage = &AuthStorage{}

var (
	ErrIsAlreadyRegistred = errors.New("user is already registred")
	ErrIsNotRegistred     = errors.New("user isn't registred")
)

type AuthStorage struct {
	// TODO: Change to db
	db map[string]*UserData
}

func New() *AuthStorage {
	return &AuthStorage{
		db: make(map[Login]*UserData),
	}
}

func (s *AuthStorage) SaveUserInfo(login Login, password Password, key Key) error {
	_, ok := s.db[login]
	if ok {
		return ErrIsAlreadyRegistred
	}

	s.db[login] = &UserData{login: login, passsword: password, key: key}

	return nil
}

func (s *AuthStorage) GetUserInfo(login string) (*UserData, error) {
	data, ok := s.db[login]
	if !ok {
		return nil, ErrIsNotRegistred
	}

	return data, nil
}
