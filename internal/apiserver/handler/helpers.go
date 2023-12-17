package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func readAuthInfoFromRequest(r *http.Request) (*AuthInfo, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read from body err=%w", err)
	}

	userData := NewAuthInfo()
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
		return "", fmt.Errorf("get Authorization cookie, err=%w", err)
	}

	return authorizationCookie.Value, nil
}

var ErrUserIsNotAuthentificated = errors.New("useris not authentificated")

func checkUserAuthorization(r *http.Request, checker AuthChecker) (string, error) {
	userKey, err := readAuthCookie(r)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", ErrUserIsNotAuthentificated
		}

		return "", fmt.Errorf("read auth cookie, err=%w", err)
	}

	login, err := checker.Check(r.Context(), userKey)
	if err != nil {
		return "", fmt.Errorf("check auth for cookie=%s, err=%w", userKey, err)
	}

	return login, nil
}

func validateOrderID(id string) bool {
	// luna validation algorithm https://ru.wikipedia.org/wiki/%D0%90%D0%BB%D0%B3%D0%BE%D1%80%D0%B8%D1%82%D0%BC_%D0%9B%D1%83%D0%BD%D0%B0

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
