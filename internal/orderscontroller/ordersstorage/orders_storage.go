package ordersstorage

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/sql"
)

var ErrLoginConflict = errors.New("order login conflict")
var ErrSavingError = errors.New("saving order error")

type Storage interface {
	HaveOrder(ctx context.Context, login string, orderId string) (bool, error)
	SaveOrder(ctx context.Context, login string, orderId string) error
}

type OrdersStorage struct {
	sqlCtrl *sql.Controller
}

func New(sqlCtrl *sql.Controller) *OrdersStorage {
	return &OrdersStorage{
		sqlCtrl: sqlCtrl,
	}
}

func (s *OrdersStorage) HaveOrder(ctx context.Context, login string, orderId string) (bool, error) {
	_, err := s.sqlCtrl.FindOrder(ctx, login, orderId)
	if err != nil {
		if errors.Is(err, sql.ErrOrderIsNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("find order login=%s order=%s, err=%w", login, orderId, err)
	}

	return true, nil
}
func (s *OrdersStorage) SaveOrder(ctx context.Context, login string, orderId string) error {
	err := s.sqlCtrl.CreateOrder(ctx, login, orderId)
	if err != nil {
		return fmt.Errorf("create order err=%w", err)
	}

	return nil
}
