package accrual

import (
	"context"
	"gophermart/internal/orderscontroller/accrual/client"
	"gophermart/internal/storage/accrualstorage"
	"gophermart/internal/storage/ordersstorage"
	"gophermart/internal/storage/sql"
	"gophermart/internal/zlog"
	"time"
)

type AccrualController struct {
	client         *client.AccrualClient
	accrualStorage accrualstorage.Storage
	ordersStorage  ordersstorage.Storage

	updaterCh chan []*sql.Order

	done chan struct{}
	wait chan struct{}
}

func StartNewController(
	ordersStorage ordersstorage.Storage,
	accrualStorage accrualstorage.Storage,
	addr string,
) *AccrualController {
	controller := &AccrualController{
		client:         client.New(addr),
		accrualStorage: accrualStorage,
		ordersStorage:  ordersStorage,

		updaterCh: make(chan []*sql.Order, 1),
		done:      make(chan struct{}),
		wait:      make(chan struct{}),
	}

	go func() {
		checkAccrualTicker := time.NewTicker(time.Second * 5)

		for {
			select {
			case <-checkAccrualTicker.C:
				if err := controller.checkAccrual(); err != nil {
					zlog.Logger.Errorf("check accrual err=%s", err)
				}
			case <-controller.done:
				close(controller.wait)

				return
			}
		}
	}()

	controller.startOrdersUpdater()

	return controller
}

func (c *AccrualController) Stop() {
	close(c.done)

	select {
	case <-c.wait:
		return
	case <-time.After(time.Second * 10):
		return
	}
}

func (c *AccrualController) checkAccrual() error {
	ctx := context.Background()

	orders, err := c.ordersStorage.GetUnexecutedOrders(ctx)
	if err != nil {
		return err
	}

	updatedOrders := c.checkOrdersStatus(ctx, orders)
	c.handleUpdatedOrders(updatedOrders)

	return nil
}

func (c *AccrualController) handleUpdatedOrders(orders []*sql.Order) {
	c.updaterCh <- orders
}

func (c *AccrualController) startOrdersUpdater() {
	go func() {
		for {
			select {
			case orders := <-c.updaterCh:
				for _, order := range orders {
					c.updateOrder(order)
				}
			case <-c.done:
				return
			}
		}
	}()
}

func (c *AccrualController) updateOrder(order *sql.Order) {
	ctx := context.Background()

	err := c.accrualStorage.UpdateAccrual(ctx, order)
	if err != nil {
		zlog.Logger.Errorf("update accrual err=%s", err)
	}
}

func (c *AccrualController) checkOrdersStatus(ctx context.Context, orders []*sql.Order) []*sql.Order {
	updatedOrders := make([]*sql.Order, 0)

	for _, order := range orders {

		accrualResponse, err := c.client.UpdateOrderStatus(ctx, order.ID)
		if err != nil {
			zlog.Logger.Errorf("accrual order=%s, err=%s", order.ID, err)
			continue
		}

		switch accrualResponse.Status {
		case client.RegistredStatus:
			if order.Status == sql.OrderStatusNew {
				order.Status = sql.OrderStatusProcessing
				updatedOrders = append(updatedOrders, order)
			}
		case client.InvalidStatus:
			order.Status = sql.OrderStatusInvalid
			updatedOrders = append(updatedOrders, order)
		case client.ProcessedStatus:
			order.Status = sql.OrderStatusProcessed
			order.Accrual = accrualResponse.Accrual
			updatedOrders = append(updatedOrders, order)
		case client.ProcessingStatus:
		default:
			zlog.Logger.Errorf("unknown accrual order status=%s", accrualResponse.Status)
		}
	}

	return updatedOrders
}
