package orderscontroller

import (
	"context"
	"errors"
	"fmt"
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
	sqlController *sql.Controller
}

func NewOrdersController(sqlController *sql.Controller) *OrdersController {
	return &OrdersController{
		sqlController: sqlController,
	}
}

func (c *OrdersController) AddOrder(ctx context.Context, login string, orderID string) (OrderStatus, error) {
	order, err := c.sqlController.FindOrder(ctx, orderID)
	if err != nil {
		if !errors.Is(err, sql.ErrOrderIsNotFound) {
			return OrderUnknownStatus, fmt.Errorf("find user=%s order=%s err=%w", login, orderID, err)
		}

		if err := c.sqlController.CreateOrder(ctx, login, orderID); err != nil {
			if errors.Is(err, sql.ErrOrderAlreadyExist) {
				return OrderAlreadyExistsStatus, nil
			}

			return OrderUnknownStatus, err
		}

		return OrderCreatedStatus, nil
	}

	if order.User != login {
		return OrderUnknownStatus, ErrOrderRegistredByOtherUser
	}

	return OrderAlreadyExistsStatus, nil
}

func (c *OrdersController) GerOrders(ctx context.Context, login string) ([]*sql.Order, error) {
	orders, err := c.sqlController.GetUserOrders(ctx, login)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, ErrOrdersListEmpty
	}

	return orders, nil
}
