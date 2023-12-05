package message

import (
	"encoding/json"
	"errors"
)

var _ Serializable = &UserData{}
var _ Desirializable = &UserData{}

var (
	ErrDesirializeData = errors.New("user data desirialization failed")
)

type UserData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewUserData() *UserData {
	return &UserData{}
}

func (m *UserData) Serialize() ([]byte, error) {
	return json.Marshal(*m)
}

func (m *UserData) Desirialize(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return errors.Join(ErrDesirializeData, err)
	}

	if m.Login == "" || m.Password == "" {
		return ErrDesirializeData
	}

	return nil
}
