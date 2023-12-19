package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/internal/sql"
	"gophermart/internal/zlog"
	"io"
	"net/http"
)

type WithdrawalsHandler struct {
	authChecker AuthChecker
	sqlCtrl     *sql.Controller
}

type WithdrawRequest struct {
	OrderID string  `json:"order"`
	Sum     float64 `json:"sum"`
}

func NewWithdrawalsHandler(authChecker AuthChecker, sqlCtrl *sql.Controller) *WithdrawalsHandler {
	return &WithdrawalsHandler{
		authChecker: authChecker,
		sqlCtrl:     sqlCtrl,
	}
}

func (h *WithdrawalsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Debugf("Withdraw handler")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := h.handle(r)
	if err != nil {
		zlog.Logger.Errorf("hanle withdraw err=%s", err)

		if errors.Is(err, ErrUserIsNotAuthentificated) {
			w.WriteHeader(http.StatusUnauthorized)
		} else if errors.Is(err, ErrBadOrderID) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else if errors.Is(err, sql.ErrNotEnoughFundsInTheAccount) {
			w.WriteHeader(http.StatusPaymentRequired)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WithdrawalsHandler) handle(r *http.Request) error {
	login, err := checkUserAuthorization(r, h.authChecker)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read from bidy err=%w", err)
	}

	req := &WithdrawRequest{}
	if err := json.Unmarshal(data, req); err != nil {
		return fmt.Errorf("unmarsgal data=%s err=%w", string(data), err)
	}

	if !validateOrderID(req.OrderID) {
		return ErrBadOrderID
	}

	if err := h.sqlCtrl.Withdraw(r.Context(), login, req.OrderID, req.Sum); err != nil {
		return fmt.Errorf("withdraw req=%v, err=%w", req, err)
	}

	return nil
}
