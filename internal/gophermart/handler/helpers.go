package handler

import (
	"fmt"
	"gophermart/internal/gophermart/userdata"
	"io"
	"net/http"
)

func readUserDataFromRequest(r *http.Request) (*userdata.AuthInfo, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read from body err=%w", err)
	}

	userData := userdata.NewUserData()
	if err := userData.Desirialize(data); err != nil {
		return nil, fmt.Errorf("user data desirialize, err=%w", err)
	}

	return userData, nil
}

func writeAuthCookie(w http.ResponseWriter, userKey string) {
	cookie := http.Cookie{Name: "Authorization", Value: string(userKey)}
	http.SetCookie(w, &cookie)
}

func readAuthCookie(r *http.Request) (string, error) {
	authorizationCookie, err := r.Cookie("Authorization")
	if err != nil {
		return "", fmt.Errorf("get Authorization cookie")
	}

	return authorizationCookie.Value, nil
}
