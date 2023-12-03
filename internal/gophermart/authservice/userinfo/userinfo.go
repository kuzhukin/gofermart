package userinfo

import "fmt"

type UserKey string

type UserInfo struct {
	Login    string
	Password string
}

func New(login string, password string) *UserInfo {
	return &UserInfo{Login: login, Password: password}
}

func (i *UserInfo) String() string {
	return fmt.Sprintf("%s-%s", i.Login, i.Password)
}
