package accrualstorage

import (
	"context"
	"fmt"
	"gophermart/internal/storage/sql"
)

type Storage interface {
	UpdateAccrual(ctx context.Context, order *sql.Order) error
}

type AccrualStorage struct {
	sqlCtrl *sql.Controller
}

func New(sqlCtrl *sql.Controller) *AccrualStorage {
	return &AccrualStorage{
		sqlCtrl: sqlCtrl,
	}
}

func (s *AccrualStorage) UpdateAccrual(ctx context.Context, order *sql.Order) error {
	if err := s.sqlCtrl.UpdateAccrual(ctx, order); err != nil {
		return fmt.Errorf("udpate accrual order=%s, user=%s, err=%w", order.ID, order.User, err)
	}

	return nil
}
