package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/internal/sql"
	"gophermart/internal/zlog"
	"net/http"
)

type BalanceWithdrawHandler struct {
	authChecker   AuthChecker
	sqlController *sql.Controller
}

func NewBalanceWithdrawHandler(authChecker AuthChecker, ctrl *sql.Controller) *BalanceWithdrawHandler {
	return &BalanceWithdrawHandler{
		authChecker:   authChecker,
		sqlController: ctrl,
	}
}

func (h *BalanceWithdrawHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := h.handle(r)
	if err != nil {
		zlog.Logger.Errorf("handle get balance, err=%s", err)

		if errors.Is(err, ErrUserIsNotAuthentificated) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	if data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		zlog.Logger.Errorf("write data, err=%s", err)
	}
}

func (h *BalanceWithdrawHandler) handle(r *http.Request) ([]byte, error) {
	login, err := checkUserAuthorization(r, h.authChecker)
	if err != nil {
		return nil, err
	}

	withdrawals, err := h.sqlController.GetUserWithdrawals(r.Context(), login)
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, nil
	}

	data, err := json.Marshal(withdrawals)
	if err != nil {
		return nil, fmt.Errorf("marshal withdrawals of user=%s, err=%w", login, err)
	}

	return data, nil
}
