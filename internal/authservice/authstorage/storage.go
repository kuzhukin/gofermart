package authstorage

import (
	"errors"
	"gophermart/internal/sql"
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

var _ Storage = &AuthStorageImpl{}

var ErrIsAlreadySaved = errors.New("is already saved")
var ErrIsNotContains = errors.New("isn't contains")

type AuthStorageImpl struct {
	sqlController *sql.Controller
}

func New(ctrl *sql.Controller) *AuthStorageImpl {
	return &AuthStorageImpl{
		sqlController: ctrl,
	}
}

func (s *AuthStorageImpl) SaveUserInfo(login Login, password Password, key Key) error {
	// s.Lock()
	// defer s.Unlock()

	// _, ok := s.login2Info[login]
	// if ok {
	// 	return ErrIsAlreadySaved
	// }

	// userInfo := &UserInfo{Login: login, Password: password, Key: key}
	// s.login2Info[login] = userInfo
	// s.key2Info[key] = userInfo

	// return nil
	return nil
}

func (s *AuthStorageImpl) GetUserInfo(login string) (*UserInfo, error) {
	// s.Lock()
	// defer s.Unlock()

	// data, ok := s.login2Info[login]
	// if !ok {
	// 	return nil, ErrIsNotContains
	// }

	// return data, nil
	return nil, nil
}

func (s *AuthStorageImpl) GetUserInfoByKey(key string) (*UserInfo, error) {
	// s.Lock()
	// defer s.Unlock()

	// data, ok := s.key2Info[key]
	// if !ok {
	// 	return nil, ErrIsNotContains
	// }

	// return data, nil
	return nil, nil
}
