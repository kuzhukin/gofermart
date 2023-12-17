package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/internal/balancecontroller"
	"gophermart/internal/zlog"
	"net/http"
)

type BalanceHandler struct {
	authChecker AuthChecker
	balanceCtrl *balancecontroller.Controller
}

func NewBalanceHandler(authChecker AuthChecker, ctrl *balancecontroller.Controller) *BalanceHandler {
	return &BalanceHandler{
		authChecker: authChecker,
		balanceCtrl: ctrl,
	}
}

func (h *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := h.handle(w, r)
	if err != nil {
		zlog.Logger.Errorf("handle get balance, err=%s", err)

		if errors.Is(err, ErrUserIsNotAuthentificated) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		zlog.Logger.Errorf("write data, err=%s", err)
	}
}

func (h *BalanceHandler) handle(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	login, err := checkUserAuthorization(r, h.authChecker)
	if err != nil {
		return nil, err
	}

	balance, err := h.balanceCtrl.GetBalnce(r.Context(), login)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(balance)
	if err != nil {
		return nil, fmt.Errorf("marshal balance=%v of user=%s, err=%w", balance, login, err)
	}

	return data, nil
}
