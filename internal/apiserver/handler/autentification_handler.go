package handler

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/zlog"
	"net/http"
)

type UserAuthorizer interface {
	Authorize(ctx context.Context, login string, password string) (string, error)
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

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	authKey, err := h.handle(r)
	if err != nil {
		zlog.Logger.Infof("Handle request was failed with err=%s", err)

		if errors.Is(err, ErrDesirializeAuthInfo) {
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
	userData, err := readAuthInfoFromRequest(r)
	if err != nil {
		return "", err
	}

	key, err := h.authorizer.Authorize(r.Context(), userData.Login, userData.Password)
	if err != nil {
		return "", fmt.Errorf("cann't authorize user with login=%s, err=%w", userData.Login, err)
	}

	return key, nil
}
