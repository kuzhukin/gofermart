package handler

import (
	"fmt"
	"gophermart/internal/gophermart/handler/message"
	"io"
	"net/http"
)

func readUserDataFromRequest(r *http.Request) (*message.UserData, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read from body err=%w", err)
	}

	userData := message.NewUserData()
	if err := userData.Desirialize(data); err != nil {
		return nil, fmt.Errorf("user data desirialize, err=%w", err)
	}

	return userData, nil
}

func writeAuthCookie(w http.ResponseWriter, userKey string) {
	cookie := http.Cookie{Name: "Authorization", Value: string(userKey)}
	http.SetCookie(w, &cookie)
}
