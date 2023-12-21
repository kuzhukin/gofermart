package handler

import (
	"encoding/json"
	"errors"
)

var (
	ErrDesirializeAuthInfo = errors.New("auth info desirialization failed")
	ErrBadAuthInfo         = errors.New("login or password isn't correct")
)

type AuthInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewAuthInfo() *AuthInfo {
	return &AuthInfo{}
}

func (m *AuthInfo) Serialize() ([]byte, error) {
	return json.Marshal(*m)
}

func (m *AuthInfo) Desirialize(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return errors.Join(ErrDesirializeAuthInfo, err)
	}

	if m.Login == "" || m.Password == "" {
		return ErrDesirializeAuthInfo
	}

	return nil
}
