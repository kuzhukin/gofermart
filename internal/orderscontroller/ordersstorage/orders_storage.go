package ordersstorage

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/sql"
)

type Storage interface {
	HaveOrder(ctx context.Context, login string, orderID string) (bool, error)
	SaveOrder(ctx context.Context, login string, orderID string) error
	GetAllOrders(ctx context.Context, login string) ([]*sql.Order, error)
}

type OrdersStorage struct {
	sqlCtrl *sql.Controller
}

func New(sqlCtrl *sql.Controller) *OrdersStorage {
	return &OrdersStorage{
		sqlCtrl: sqlCtrl,
	}
}

func (s *OrdersStorage) HaveOrder(ctx context.Context, login string, orderID string) (bool, error) {
	_, err := s.sqlCtrl.FindOrder(ctx, login, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrOrderIsNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("find order login=%s order=%s, err=%w", login, orderID, err)
	}

	return true, nil
}

func (s *OrdersStorage) SaveOrder(ctx context.Context, login string, orderID string) error {
	err := s.sqlCtrl.CreateOrder(ctx, login, orderID)
	if err != nil {
		return fmt.Errorf("create order err=%w", err)
	}

	return nil
}

func (s *OrdersStorage) GetAllOrders(ctx context.Context, login string) ([]*sql.Order, error) {
	orders, err := s.sqlCtrl.GetAllOrders(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get orders for user=%s, err=%w", login, err)
	}

	return orders, nil
}
