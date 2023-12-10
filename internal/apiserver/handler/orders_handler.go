package handler

import (
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
	Check(userKey string) (string, error)
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

	login, err := h.authChecker.Check(userKey)
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
		} else if errors.Is(err, ErrBadOrderId) {
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
	orderId, err := getOrderId(r)
	if err != nil {
		return orderscontroller.OrderUnknownStatus, fmt.Errorf("get order id, err=%w", err)
	}

	status, err := h.orderscontroller.AddOrder(r.Context(), login, orderId)
	if err != nil {
		return orderscontroller.OrderUnknownStatus, fmt.Errorf("add order, err=%w", err)
	}

	return status, nil
}

var ErrBadOrderId = errors.New("bad order id")
var ErrBadRequestFormat = errors.New("can't read order id from body")

func getOrderId(r *http.Request) (string, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.Join(ErrBadRequestFormat, err)
	}

	id, err := parseOrderId(data)
	if err != nil {
		return "", fmt.Errorf("parse orderId=%s, err=%w", string(data), err)
	}

	return id, nil
}

func parseOrderId(data []byte) (string, error) {
	orderId := string(data)

	if len(orderId) == 0 {
		return "", ErrBadRequestFormat
	}

	// TODO: add validation bu luna's algorithm https://ru.wikipedia.org/wiki/%D0%90%D0%BB%D0%B3%D0%BE%D1%80%D0%B8%D1%82%D0%BC_%D0%9B%D1%83%D0%BD%D0%B0

	for _, c := range orderId {
		if !unicode.IsDigit(c) {
			return "", ErrBadOrderId
		}
	}

	if ok := validateOrderId(orderId); !ok {
		return "", ErrBadOrderId
	}

	return orderId, nil
}

func (h *OrdersHandler) serveGetOrderList(w http.ResponseWriter, r *http.Request, login string) {

}
