package orderscontroller

import (
	"context"
	"errors"
	"gophermart/internal/orderscontroller/ordersstorage"
	"gophermart/internal/sql"
)

type OrderStatus int

const (
	OrderUnknownStatus       OrderStatus = 0
	OrderAlreadyExistsStatus OrderStatus = 1
	OrderCreatedStatus       OrderStatus = 2
)

var ErrOrderRegistredByOtherUser = errors.New("order was registred by other user")
var ErrOrdersListEmpty = errors.New("orders list empty")

type OrdersController struct {
	storage ordersstorage.Storage
}

func NewOrdersController(storage ordersstorage.Storage) *OrdersController {
	return &OrdersController{
		storage: storage,
	}
}

func (c *OrdersController) AddOrder(ctx context.Context, login string, orderID string) (OrderStatus, error) {
	if err := c.storage.SaveOrder(ctx, login, orderID); err != nil {
		if errors.Is(err, sql.ErrOrderAlreadyExist) {
			return OrderAlreadyExistsStatus, nil
		}

		return OrderUnknownStatus, err
	}

	return OrderCreatedStatus, nil
}

func (c *OrdersController) GerOrders(ctx context.Context, login string) ([]*sql.Order, error) {
	orders, err := c.storage.GetAllOrders(ctx, login)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, ErrOrdersListEmpty
	}

	return orders, nil
}
