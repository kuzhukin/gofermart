package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/internal/sql"
	"gophermart/internal/zlog"
	"net/http"
)

type BalanceHandler struct {
	authChecker   AuthChecker
	sqlController *sql.Controller
}

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func NewBalanceHandler(authChecker AuthChecker, ctrl *sql.Controller) *BalanceHandler {
	return &BalanceHandler{
		authChecker:   authChecker,
		sqlController: ctrl,
	}
}

func (h *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		zlog.Logger.Errorf("write data, err=%s", err)
	}
}

func (h *BalanceHandler) handle(r *http.Request) ([]byte, error) {
	login, err := checkUserAuthorization(r, h.authChecker)
	if err != nil {
		return nil, err
	}

	userStatistic, err := h.sqlController.GetUserStatistic(r.Context(), login)
	if err != nil {
		return nil, err
	}

	response := &BalanceResponse{Current: userStatistic.Balance, Withdrawn: userStatistic.WithdrawalsTotalSum}

	data, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("marshal balance=%v of user=%s, err=%w", response, login, err)
	}

	return data, nil
}
