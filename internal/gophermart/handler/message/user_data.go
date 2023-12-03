package message

import (
	"encoding/json"
)

var _ Serializable = &UserData{}
var _ Desirializable = &UserData{}

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
	return json.Unmarshal(data, m)
}
