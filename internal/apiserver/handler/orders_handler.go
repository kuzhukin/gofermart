package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/internal/orderscontroller"
	"gophermart/internal/zlog"
	"io"
	"net/http"
	"unicode"
)

var ErrUserIsNotAuthentificated = errors.New("user is not authentificated")
var ErrUnsuportedMethod = errors.New("unsuported method")

type AuthChecker interface {
	Check(ctx context.Context, userKey string) (string, error)
}

type OrdersHandler struct {
	authChecker      AuthChecker
	orderscontroller *orderscontroller.OrdersController
}

func NewOrdersHandler(authChecker AuthChecker, ordersController *orderscontroller.OrdersController) *OrdersHandler {
	return &OrdersHandler{
		authChecker:      authChecker,
		orderscontroller: ordersController,
	}
}

func (h *OrdersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Debugf("Orders handler")

	handler, err := h.selectHandler(r.Method)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login, err := h.checkUserAuthorization(r)
	if err != nil {
		if errors.Is(err, ErrUserIsNotAuthentificated) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}

	handler(w, r, login)
}

func (h *OrdersHandler) checkUserAuthorization(r *http.Request) (string, error) {
	userKey, err := readAuthCookie(r)
	if err != nil {
		return "", fmt.Errorf("read auth cookie, err=%w", err)
	}

	login, err := h.authChecker.Check(r.Context(), userKey)
	if err != nil {
		return "", fmt.Errorf("check auth for cookie=%s, err=%w", userKey, err)
	}

	return login, nil
}

func (h *OrdersHandler) selectHandler(method string) (func(w http.ResponseWriter, r *http.Request, login string), error) {
	switch method {
	case http.MethodPost:
		return h.serveLoadNewOrder, nil
	case http.MethodGet:
		return h.serveGetOrderList, nil
	default:
		return nil, ErrUnsuportedMethod
	}
}

func (h *OrdersHandler) serveLoadNewOrder(w http.ResponseWriter, r *http.Request, login string) {
	orderStatus, err := h.loadNewOrder(w, r, login)
	if err != nil {
		zlog.Logger.Infof("Load new order user=%s err=%s", login, err)

		if errors.Is(err, ErrBadRequestFormat) {
			w.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(err, ErrBadOrderID) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else if errors.Is(err, orderscontroller.ErrOrderRegistredByOtherUser) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	if orderStatus == orderscontroller.OrderCreatedStatus {
		w.WriteHeader(http.StatusAccepted)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrdersHandler) loadNewOrder(w http.ResponseWriter, r *http.Request, login string) (orderscontroller.OrderStatus, error) {
	orderID, err := getOrderID(r)
	if err != nil {
		return orderscontroller.OrderUnknownStatus, fmt.Errorf("get order id, err=%w", err)
	}

	status, err := h.orderscontroller.AddOrder(r.Context(), login, orderID)
	if err != nil {
		return orderscontroller.OrderUnknownStatus, fmt.Errorf("add order, err=%w", err)
	}

	return status, nil
}

var ErrBadOrderID = errors.New("bad order id")
var ErrBadRequestFormat = errors.New("can't read order id from body")

func getOrderID(r *http.Request) (string, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.Join(ErrBadRequestFormat, err)
	}

	id, err := parseOrderID(data)
	if err != nil {
		return "", fmt.Errorf("parse orderID=%s, err=%w", string(data), err)
	}

	return id, nil
}

func parseOrderID(data []byte) (string, error) {
	orderID := string(data)

	if len(orderID) == 0 {
		return "", ErrBadRequestFormat
	}

	for _, c := range orderID {
		if !unicode.IsDigit(c) {
			return "", ErrBadOrderID
		}
	}

	if ok := validateOrderID(orderID); !ok {
		return "", ErrBadOrderID
	}

	return orderID, nil
}

func (h *OrdersHandler) serveGetOrderList(w http.ResponseWriter, r *http.Request, login string) {
	data, err := h.getUserOrders(r, login)
	if err != nil {
		zlog.Logger.Errorf("Get user=%s orders err=%s", login, err)

		if errors.Is(err, orderscontroller.ErrOrdersListEmpty) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *OrdersHandler) getUserOrders(r *http.Request, login string) ([]byte, error) {
	orders, err := h.orderscontroller.GerOrders(r.Context(), login)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}

	return data, nil
}
