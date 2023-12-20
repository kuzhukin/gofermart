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

	login, err := checkUserAuthorization(r, h.authChecker)
	if err != nil {
		zlog.Logger.Errorf("Check user auth err=%s", err)

		if errors.Is(err, ErrUserIsNotAuthentificated) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	handler, err := h.selectHandler(r.Method)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	handler(w, r, login)
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

	zlog.Logger.Debugf("GET user orders %s", string(data))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		zlog.Logger.Errorf("err write orders, err=%s", err)
	}
}

type OrderResponse struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

func (h *OrdersHandler) getUserOrders(r *http.Request, login string) ([]byte, error) {
	orders, err := h.orderscontroller.GerOrders(r.Context(), login)
	if err != nil {
		return nil, err
	}

	// reformatting
	responses := make([]*OrderResponse, 0, len(orders))
	for _, order := range orders {
		resp := &OrderResponse{
			Number:     order.ID,
			Status:     string(order.Status),
			Accrual:    float64(order.Accrual),
			UploadedAt: order.UpdaloadTime,
		}
		responses = append(responses, resp)
	}

	data, err := json.Marshal(responses)
	if err != nil {
		return nil, err
	}

	return data, nil
}
