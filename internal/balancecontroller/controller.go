package balancecontroller

import (
	"context"
	"fmt"
	"gophermart/internal/storage/userstorage"
)

type Controller struct {
	userStorage userstorage.Storage
}

func New(storage userstorage.Storage) *Controller {
	return &Controller{
		userStorage: storage,
	}
}

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (c *Controller) GetBalnce(ctx context.Context, login string) (*BalanceResponse, error) {
	user, err := c.userStorage.GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get user=%s from storage err=%w", login, err)
	}

	// TODO: add getting withdrawn

	return &BalanceResponse{Current: float64(user.Balance)}, nil
}

func (c *Controller) Withdraw(ctx context.Context, login string, orderID string, amount float64) error {
	return nil
}
