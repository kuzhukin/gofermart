package handler

import (
	"fmt"
	"gophermart/internal/userdata"
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

func validateOrderId(id string) bool {
	sum := 0
	parity := len(id) % 2

	for i := 0; i < len(id); i++ {
		digit := int(id[i] - '0')

		if (i % 2) == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return (sum % 10) == 0
}
