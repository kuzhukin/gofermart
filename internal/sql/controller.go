package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gophermart/internal/zlog"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	createTablesTimeout = time.Second * 10
	queryTimeout        = time.Second * 1
	execTimeout         = time.Second * 1
)

type Controller struct {
	db *sql.DB
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID     string
	User   string
	Status OrderStatus
}

type User struct {
	Login     string
	AuthToken string
	Balance   uint64
}

func StartNewController(dataSourceName string) (*Controller, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("sql open db=%s err=%w", dataSourceName, err)
	}

	ctrl := &Controller{db: db}
	if err := ctrl.init(); err != nil {
		return nil, fmt.Errorf("init, err=%w", err)
	}

	return ctrl, nil
}

func (c *Controller) Stop() error {
	if err := c.db.Close(); err != nil {
		return fmt.Errorf("db close err=%w", err)
	}

	return nil
}

func (c *Controller) init() error {
	createTableQueries := []string{
		createUsersTableQuery,
		createOrdersTableQuery,
		createWithdrawalsTableQuery,
	}

	for _, q := range createTableQueries {
		if err := c.exec(q); err != nil {
			return fmt.Errorf("exec query=%s, err=%w", q, err)
		}
	}

	return nil
}

func (c *Controller) exec(query string) error {
	ctx, cancel := context.WithTimeout(context.Background(), createTablesTimeout)
	defer cancel()

	_, err := c.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("exec create table query, err=%w", err)
	}

	return nil
}

// ----------------------------------------------------------------------------------------------
// ---------------------------------------- User Methods ----------------------------------------
// ----------------------------------------------------------------------------------------------

var ErrUserIsNotFound = errors.New("user isn't found")

func (c *Controller) CreateUser(ctx context.Context, login string, token string) error {
	queryFunc := c.makeExecFunc(ctx, createUserQuery, []interface{}{login, token})

	_, err := doQuery(queryFunc)
	if err != nil {
		return fmt.Errorf("exec create user err=%w", err)
	}

	return nil
}

func (c *Controller) FindUser(ctx context.Context, login string) (*User, error) {
	queryFunc := c.makeQueryFunc(ctx, getUser, []interface{}{login})

	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("do find user query err=%w", err)
	}

	user := &User{}

	if rows.Next() {
		if err := rows.Scan(&user.Login, &user.AuthToken, &user.Balance); err != nil {
			return nil, fmt.Errorf("rows scan to user, err=%w", err)
		}

		return user, nil
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Errorf("rows close err=%s", err)
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrUserIsNotFound
}

func (c *Controller) FindUserByToken(ctx context.Context, token string) (*User, error) {
	queryFunc := c.makeQueryFunc(ctx, getUserByToken, []interface{}{token})

	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("do find user by token query err=%w", err)
	}

	user := &User{}

	if rows.Next() {
		if err := rows.Scan(&user.Login, &user.AuthToken, &user.Balance); err != nil {
			return nil, fmt.Errorf("rows scan to user, err=%w", err)
		}

		return user, nil
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Errorf("rows close err=%s", err)
		}
	}()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrUserIsNotFound
}

// func (c *Controller) Withdraw(ctx cont)

// -----------------------------------------------------------------------------------------------
// ------------------------------------- Orders handling API -------------------------------------
// -----------------------------------------------------------------------------------------------

var ErrOrderIsNotFound = errors.New("order isn't found")

func (c *Controller) FindOrder(ctx context.Context, login string, orderID string) (*Order, error) {
	queryFunc := c.makeQueryFunc(ctx, getOrderQuery, []interface{}{orderID, login})

	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("do query, err=%w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Errorf("rows close err=%s", err)
		}
	}()

	order := &Order{}

	if rows.Next() {
		if err := rows.Scan(&order.ID, &order.Status, order.User); err != nil {
			return nil, fmt.Errorf("rows scan to order err=%w", err)
		}

		return order, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrOrderIsNotFound
}

func (c *Controller) CreateOrder(ctx context.Context, login string, orderID string) error {
	execFunc := c.makeExecFunc(ctx, createOrderQuery, []interface{}{orderID, login})

	_, err := doQuery(execFunc)
	if err != nil {
		return fmt.Errorf("create user=%s, order=%s, err=%w", login, orderID, err)
	}

	return nil
}

// ----------------------------------------------------------------------------------------------
// -------------------------------------- Internal Methods --------------------------------------
// ----------------------------------------------------------------------------------------------

func (c *Controller) makeExecFunc(ctx context.Context, query string, args []interface{}) func() (*sql.Result, error) {
	return func() (r *sql.Result, err error) {
		ctx, cancel := context.WithTimeout(ctx, execTimeout)
		defer cancel()

		res, err := c.db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("exec query=%v with args=%v err=%w", query, args, err)
		}

		return &res, nil
	}
}

func (c *Controller) makeQueryFunc(ctx context.Context, query string, args []interface{}) func() (*sql.Rows, error) {
	return func() (*sql.Rows, error) {
		ctx, cancel := context.WithTimeout(ctx, queryTimeout)
		defer cancel()

		rows, err := c.db.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("query metric, err=%w", err)
		}

		return rows, nil
	}
}

var tryingIntervals = []time.Duration{
	time.Millisecond * 100,
	time.Millisecond * 300,
	time.Millisecond * 500,
}

func doQuery[T any](queryFunc func() (*T, error)) (*T, error) {
	var commonErr error
	max := len(tryingIntervals)

	for trying := 0; trying <= max; trying++ {
		rows, err := queryFunc()
		if err != nil {
			commonErr = errors.Join(commonErr, err)

			if trying < max && isRetriableError(err) {
				time.Sleep(tryingIntervals[trying])
				continue
			}

			return nil, commonErr
		}

		return rows, nil
	}

	return nil, commonErr
}

func isRetriableError(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && (pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) || pgerrcode.IsConnectionException(pgErr.Code))
}
