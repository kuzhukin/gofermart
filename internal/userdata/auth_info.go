package userdata

import (
	"encoding/json"
	"errors"
)

var (
	ErrDesirializeUserData = errors.New("auth info desirialization failed")
	ErrBadUserData         = errors.New("login or password isn't correct")
)

type AuthInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewUserData() *AuthInfo {
	return &AuthInfo{}
}

func (m *AuthInfo) Serialize() ([]byte, error) {
	return json.Marshal(*m)
}

func (m *AuthInfo) Desirialize(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return errors.Join(ErrDesirializeUserData, err)
	}

	if m.Login == "" || m.Password == "" {
		return ErrDesirializeUserData
	}

	return nil
}
