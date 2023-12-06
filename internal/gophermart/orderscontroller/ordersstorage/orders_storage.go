package ordersstorage

import "errors"

var ErrLoginConflict = errors.New("order login conflict")
var ErrSavingError = errors.New("saving order error")

type Storage interface {
	HaveOrder(login string, orderId string) (bool, error)
	SaveOrder(login string, orderId string) error
}

type OrdersStorage struct {
}

func New() *OrdersStorage {
	return &OrdersStorage{}
}

func (s *OrdersStorage) HaveOrder(login string, orderId string) (bool, error) {
	return false, nil
}
func (s *OrdersStorage) SaveOrder(login string, orderId string) error {
	return nil
}
