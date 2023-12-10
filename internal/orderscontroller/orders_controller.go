package orderscontroller

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/orderscontroller/ordersstorage"
)

type OrderStatus int

const (
	OrderUnknownStatus       OrderStatus = 0
	OrderAlreadyExistsStatus OrderStatus = 1
	OrderCreatedStatus       OrderStatus = 2
)

var ErrOrderRegistredByOtherUser = errors.New("order was registred by other user")

type OrdersController struct {
	storage ordersstorage.Storage
}

func NewOrdersController(storage ordersstorage.Storage) *OrdersController {
	return &OrdersController{
		storage: storage,
	}
}

func (c *OrdersController) AddOrder(ctx context.Context, login string, orderID string) (OrderStatus, error) {
	ok, err := c.storage.HaveOrder(ctx, login, orderID)
	if err != nil {
		if errors.Is(err, ordersstorage.ErrLoginConflict) {
			return OrderUnknownStatus, ErrOrderRegistredByOtherUser
		}

		return OrderUnknownStatus, fmt.Errorf("have order, err=%w", err)
	}

	if ok {
		return OrderAlreadyExistsStatus, nil
	}

	if err := c.storage.SaveOrder(ctx, login, orderID); err != nil {
		return OrderUnknownStatus, fmt.Errorf("save order, err=%w", err)
	}

	return OrderAlreadyExistsStatus, nil
}
