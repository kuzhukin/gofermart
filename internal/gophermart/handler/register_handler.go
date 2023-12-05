package handler

import (
	"errors"
	"fmt"
	"gophermart/internal/gophermart/handler/message"
	"gophermart/internal/gophermart/zlog"
	"net/http"
)

type UserKey = string

type UserRegistrator interface {
	Register(login string, password string) (UserKey, error)
}

var (
	ErrIsAlreadyRegistred = errors.New("user is already registred")
)

type RegistrationHandler struct {
	registrator UserRegistrator
}

func NewRegistrationHandler(registrator UserRegistrator) *RegistrationHandler {
	return &RegistrationHandler{
		registrator: registrator,
	}
}

func (h *RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Debugf("Register handler")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	authKey, err := h.handle(r)
	if err != nil {
		zlog.Logger.Infof("Handle request was failed with err=%s", err)

		if errors.Is(err, message.ErrDesirializeData) {
			w.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(err, ErrIsAlreadyRegistred) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	writeAuthCookie(w, authKey)
	w.WriteHeader(http.StatusOK)
}

func (h *RegistrationHandler) handle(r *http.Request) (UserKey, error) {
	userData, err := readUserDataFromRequest(r)
	if err != nil {
		return "", err
	}

	userKey, err := h.registrator.Register(userData.Login, userData.Password)
	if err != nil {
		return "", fmt.Errorf("registration failed, err=%w", err)
	}

	return userKey, nil
}
