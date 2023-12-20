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

	createUserTimeout = time.Second * 1
	getUserTimeout    = time.Second * 1
	withdrawTimeout   = time.Second * 3

	getOrderTimeout           = time.Second * 1
	getAllOrdersTimeout       = time.Second * 10
	createOrderTimeout        = time.Second * 1
	updateOrderAccrualTimeout = time.Second * 2
)

type query struct {
	request string
	args    []interface{}
}

type Controller struct {
	db     *sql.DB
	dbPath string
}

func StartNewController(dataSourceName string) (*Controller, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("sql open db=%s err=%w", dataSourceName, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db connection err=%w", err)
	}

	ctrl := &Controller{db: db, dbPath: dataSourceName}
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
	queryFunc := c.makeExecFunc(ctx, prepareCreateUserQuery(login, token))

	_, err := doQuery(queryFunc)
	if err != nil {
		return fmt.Errorf("exec create user err=%w", err)
	}

	return nil
}

func (c *Controller) FindUser(ctx context.Context, login string) (*User, error) {
	queryFunc := c.makeQueryFunc(ctx, prepareGetUserQuery(login), getUserTimeout)

	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("do find user query err=%w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Errorf("rows close err=%s", err)
		}
	}()

	user := &User{}
	if rows.Next() {
		if err := rows.Scan(&user.Login, &user.AuthToken, &user.Balance); err != nil {
			return nil, fmt.Errorf("rows scan to user, err=%w", err)
		}

		return user, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrUserIsNotFound
}

func (c *Controller) FindUserByToken(ctx context.Context, token string) (*User, error) {
	queryFunc := c.makeQueryFunc(ctx, prepareGetUserByTokenQuery(token), getUserTimeout)

	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("do find user by token query err=%w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Errorf("rows close err=%s", err)
		}
	}()

	if !rows.Next() {
		return nil, ErrUserIsNotFound
	}

	user := &User{}
	if err := user.scan(rows); err != nil {
		return nil, fmt.Errorf("rows scan to user, err=%w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return user, nil

}

var ErrNotEnoughFundsInTheAccount = errors.New("there are not enough funds in the account")

func (c *Controller) Withdraw(ctx context.Context, login string, orderID string, amount float64) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx err=%w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	user, err := doTransactionQuery(ctx, tx, prepareGetUserQuery(login), scanUserFromRows)
	if err != nil {
		return err
	}

	if user.Balance < amount {
		return ErrNotEnoughFundsInTheAccount
	}

	decreaseUserBalanceQury := prepareDecreaseUserBalanceQuery(login, amount)

	_, err = tx.ExecContext(ctx, decreaseUserBalanceQury.request, decreaseUserBalanceQury.args...)
	if err != nil {
		return fmt.Errorf("decrese user=%s balance=%.4f on amount=%.4f err=%w", user.Login, user.Balance, amount, err)
	}

	addWitdhrawalsQuery := prepareAddWithdrawalsQuery(orderID, login, amount)

	_, err = tx.ExecContext(ctx, addWitdhrawalsQuery.request, addWitdhrawalsQuery.args...)
	if err != nil {
		return fmt.Errorf("add withdrawals query orderID=%s login=%s amount=%.4f err=%w", orderID, login, amount, err)
	}

	return tx.Commit()
}

type UserStatistic struct {
	Balance             float64
	WithdrawalsTotalSum float64
}

func (c *Controller) GetUserStatistic(ctx context.Context, login string) (*UserStatistic, error) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx err=%w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	getUserQuery := prepareGetUserQuery(login)
	withdrawaslSumQuery := prepareWithdrawalsSumQuery(login)

	user, err := doTransactionQuery(ctx, tx, getUserQuery, scanUserFromRows)
	if err != nil {
		return nil, err
	}

	withdrawalsSum, err := doTransactionQuery(ctx, tx, withdrawaslSumQuery, scanWithdrawalsSumFromRows)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		zlog.Logger.Errorf("commit tx err=%s", err)
	}

	return &UserStatistic{Balance: user.Balance, WithdrawalsTotalSum: withdrawalsSum}, nil
}

func (c *Controller) GetUserWithdrawals(ctx context.Context, user string) ([]*UserWithdrawRecord, error) {
	queryFunc := c.makeQueryFunc(ctx, prepareGetAllUserWithdrawals(user), time.Second*5)
	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("do get all withdrawals of user=%s query err=%w", user, err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Errorf("rows close err=%s", err)
		}
	}()

	list := make([]*UserWithdrawRecord, 0)
	for rows.Next() {
		wr := &UserWithdrawRecord{}
		if err := wr.scan(rows); err != nil {
			return nil, fmt.Errorf("scan withdraw err=%w", err)
		}

		list = append(list, wr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err=%w", err)
	}

	return list, nil
}

// -----------------------------------------------------------------------------------------------
// ------------------------------------- Orders handling API -------------------------------------
// -----------------------------------------------------------------------------------------------

const TimestampFormat = "2006-01-02T15:04:05"

var ErrOrderIsNotFound = errors.New("order isn't found")
var ErrOrderAlreadyExist = errors.New("order already exist")

func (c *Controller) FindOrder(ctx context.Context, orderID string) (*Order, error) {
	queryFunc := c.makeQueryFunc(ctx, prepareGetOrderQuery(orderID), getOrderTimeout)

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

	if !rows.Next() {
		return nil, ErrOrderIsNotFound
	}

	order := &Order{}
	if err := order.scan(rows); err != nil {
		return nil, fmt.Errorf("rows scan to order err=%w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return order, nil
}

func (c *Controller) CreateOrder(ctx context.Context, login string, orderID string) error {
	execFunc := c.makeExecFunc(ctx, prepareCreateOrderQuery(orderID, login))

	_, err := doQuery(execFunc)
	if err != nil {
		if isNotUniqueError(err) {
			return ErrOrderAlreadyExist
		}

		return fmt.Errorf("create user=%s, order=%s, err=%w", login, orderID, err)
	}

	return nil
}

func (c *Controller) GetUserOrders(ctx context.Context, login string) ([]*Order, error) {
	queryFunc := c.makeQueryFunc(ctx, prepareGetAllOrdersQuery(login), getAllOrdersTimeout)

	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("get all orders query, err=%w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Errorf("rows close err=%s", err)
		}
	}()

	orders := make([]*Order, 0)
	for rows.Next() {
		order := &Order{}
		if err := order.scan(rows); err != nil {
			return nil, fmt.Errorf("scan rows, err=%w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (c *Controller) GetUnexecutedOrders(ctx context.Context) ([]*Order, error) {
	queryFunc := c.makeQueryFunc(ctx, prepareGetUnexecutedOrdersQuery(), getAllOrdersTimeout)

	rows, err := doQuery(queryFunc)
	if err != nil {
		return nil, fmt.Errorf("get unexecuted orders query, err=%w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Errorf("rows close err=%s", err)
		}
	}()

	orders := make([]*Order, 0)
	for rows.Next() {
		order := &Order{}
		err := order.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan rows, err=%w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (c *Controller) UpdateAccrual(ctx context.Context, order *Order) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx err=%w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, updateOrderAccrualQuery, order.Status, order.Accrual, order.ID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, increaseUserBalanceQuery, order.Accrual, order.User); err != nil {
		return err
	}

	return tx.Commit()
}

// ----------------------------------------------------------------------------------------------
// -------------------------------------- Internal Methods --------------------------------------
// ----------------------------------------------------------------------------------------------

func (c *Controller) makeExecFunc(ctx context.Context, query *query) func() (*sql.Result, error) {
	return func() (r *sql.Result, err error) {
		res, err := c.db.ExecContext(ctx, query.request, query.args...)
		if err != nil {
			return nil, fmt.Errorf("exec query=%v err=%w", query, err)
		}

		return &res, nil
	}
}

func (c *Controller) makeQueryFunc(ctx context.Context, query *query, _ time.Duration) func() (*sql.Rows, error) {
	return func() (*sql.Rows, error) {
		rows, err := c.db.QueryContext(ctx, query.request, query.args...)
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

	return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}

func isNotUniqueError(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}

func doTransactionQuery[T any](ctx context.Context, tx *sql.Tx, query *query, parseFun func(rows *sql.Rows) (T, error)) (T, error) {
	var result T

	rows, err := tx.QueryContext(ctx, query.request, query.args...)
	if err != nil {
		return result, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Errorf("rows close err=%s", err)
		}
	}()

	result, err = parseFun(rows)
	if err != nil {
		return result, err
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}
