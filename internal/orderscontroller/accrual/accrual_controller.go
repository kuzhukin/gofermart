package accrual

import (
	"context"
	"gophermart/internal/orderscontroller/accrual/client"
	"gophermart/internal/sql"
	"gophermart/internal/zlog"
	"time"
)

const pollingInterval = time.Second * 1
const shutdownTimeout = time.Second * 10

type AccrualController struct {
	client        *client.AccrualClient
	sqlController *sql.Controller

	updaterCh chan []*sql.Order

	done chan struct{}
	wait chan struct{}
}

func StartNewController(
	sqlController *sql.Controller,
	addr string,
) *AccrualController {
	controller := &AccrualController{
		client:        client.New(addr),
		sqlController: sqlController,

		updaterCh: make(chan []*sql.Order, 1),
		done:      make(chan struct{}),
		wait:      make(chan struct{}),
	}

	go func() {
		checkAccrualTicker := time.NewTicker(pollingInterval)

		for {
			select {
			case <-checkAccrualTicker.C:
				zlog.Logger.Debugf("Start check accrual")

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
	case <-time.After(shutdownTimeout):
		return
	}
}

func (c *AccrualController) checkAccrual() error {
	ctx := context.Background()

	orders, err := c.sqlController.GetUnexecutedOrders(ctx)
	if err != nil {
		return err
	}

	zlog.Logger.Debugf("accrual controller has %d unexecuted orders", len(orders))

	updatedOrders := c.checkOrdersStatus(ctx, orders)

	zlog.Logger.Debugf("accrual controller has %d updated orders", len(updatedOrders))

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
				zlog.Logger.Debugf("START UPDATE ORDERS count=%d", len(orders))
				for _, order := range orders {
					zlog.Logger.Debugf("UPDATE ORDERS %+v", order)
					c.updateOrder(order)
					time.Sleep(time.Millisecond * 100)
				}
			case <-c.done:
				return
			}
		}
	}()
}

func (c *AccrualController) updateOrder(order *sql.Order) {
	ctx := context.Background()

	zlog.Logger.Debugf("write new order state to db order=%+v", order)

	err := c.sqlController.UpdateAccrual(ctx, order)
	if err != nil {
		zlog.Logger.Errorf("update accrual err=%s", err)
	}
}

func (c *AccrualController) checkOrdersStatus(ctx context.Context, orders []*sql.Order) []*sql.Order {
	updatedOrders := make([]*sql.Order, 0)

	for _, order := range orders {

		accrualResponse, err := c.client.UpdateOrderStatus(ctx, order.ID)
		if err != nil {
			zlog.Logger.Infof("accrual order=%s, err=%s", order.ID, err)
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
