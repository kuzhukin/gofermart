package balancecontroller

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/sql"
)

var ErrNotEnoughFundsInTheAccount = errors.New("there are not enough funds in the account")

type Controller struct {
	sqlController *sql.Controller
}

func New(sqlController *sql.Controller) *Controller {
	return &Controller{
		sqlController: sqlController,
	}
}

type BalanceResponse struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

func (c *Controller) GetBalnce(ctx context.Context, login string) (*BalanceResponse, error) {
	user, err := c.sqlController.FindUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get user=%s from storage err=%w", login, err)
	}

	// TODO: add getting withdrawn

	return &BalanceResponse{Current: user.Balance}, nil
}

func (c *Controller) Withdraw(ctx context.Context, login string, orderID string, amount float32) error {
	return ErrNotEnoughFundsInTheAccount
}
