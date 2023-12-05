package handler

import (
	"errors"
	"fmt"
	"gophermart/internal/gophermart/authservice/storage"
	"gophermart/internal/gophermart/handler/message"
	"gophermart/internal/gophermart/zlog"
	"io"
	"net/http"
)

var _ http.Handler = &RegisterHandler{}

type UserKey = string

type UserRegistrator interface {
	Register(login string, password string) (UserKey, error)
}

type RegisterHandler struct {
	registrator UserRegistrator
}

func NewRegisterHandler(registrator UserRegistrator) *RegisterHandler {
	return &RegisterHandler{
		registrator: registrator,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Debugf("Register handler")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userKey, err := h.handle(r)
	if err != nil {
		zlog.Logger.Infof("Handle request was failed with err=%s", err)

		if errors.Is(err, message.ErrDesirializeData) {
			w.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(err, storage.ErrIsAlreadyRegistred) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	cookie := http.Cookie{Name: "Authorization", Value: string(userKey)}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *RegisterHandler) handle(r *http.Request) (UserKey, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("read from body err=%w", err)
	}

	msg := message.NewUserData()
	if err := msg.Desirialize(data); err != nil {
		return "", fmt.Errorf("user data desirialize, err=%w", err)
	}

	userKey, err := h.registrator.Register(msg.Login, msg.Password)
	if err != nil {
		return "", fmt.Errorf("registration failed, err=%w", err)
	}

	return userKey, nil
}
