package handler

import (
	"errors"
	"fmt"
	"gophermart/internal/gophermart/handler/message"
	"gophermart/internal/gophermart/zlog"
	"net/http"
)

type UserAuthorizer interface {
	Authorize(login string, password string) (string, error)
}

var (
	ErrIsNotAutorized = errors.New("user isn't autorized")
)

type AutentifiactionHandler struct {
	authorizer UserAuthorizer
}

func NewAutentifiactionHandler(authorizer UserAuthorizer) *AutentifiactionHandler {
	return &AutentifiactionHandler{
		authorizer: authorizer,
	}
}

func (h *AutentifiactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Debugf("auth handler")

	authKey, err := h.handle(r)
	if err != nil {
		if errors.Is(err, message.ErrDesirializeData) {
			w.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(err, ErrIsNotAutorized) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	writeAuthCookie(w, authKey)
	w.WriteHeader(http.StatusOK)
}

func (h *AutentifiactionHandler) handle(r *http.Request) (string, error) {
	userData, err := readUserDataFromRequest(r)
	if err != nil {
		return "", err
	}

	key, err := h.authorizer.Authorize(userData.Login, userData.Password)
	if err != nil {
		return "", fmt.Errorf("authorize login=%s", userData.Login)
	}

	return key, nil
}
