package sql

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	// BONUSES WITHDRAWALS TABLE
	createWithdrawalsTableQuery = `CREATE TABLE IF NOT EXISTS withdrawals (
		"order"			text				NOT NULL,
		"user"			text				NOT NULL,
		"sum"			double precision	NOT NULL,
		"processed_at"	text				NOT NULL,
		PRIMARY KEY ( "order" )
	);`

	addWithdrawals             = `INSERT INTO withdrawals ("order", "user", "sum", "processed_at") VALUES ($1, $2, $3, $4);`
	getUserWithdrawalsTotalSum = `SELECT SUM ("sum") FROM withdrawals WHERE "user" = $1;`
	getUserWithdrawals         = `SELECT "order", "sum", "processed_at" FROM withdrawals WHERE "user" = $1 ORDER BY "processed_at";`
)

type UserWithdrawRecord struct {
	OrderID     string  `json:"order"`
	Accrual     float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (r *UserWithdrawRecord) scan(rows *sql.Rows) error {
	err := rows.Scan(&r.OrderID, &r.Accrual, &r.ProcessedAt)
	if err != nil {
		return fmt.Errorf("withdraw record scan err=%w", err)
	}

	return nil
}

func scanWithdrawalsSumFromRows(rows *sql.Rows) (float64, error) {
	if !rows.Next() {
		return 0, ErrEmptyScannerResult
	}

	var result interface{}
	if err := rows.Scan(&result); err != nil {
		return 0, fmt.Errorf("scan withdraws sum err=%w", err)
	}

	var withdrawsSum float64
	switch v := result.(type) {
	case float32:
		withdrawsSum = float64(v)
	case float64:
		withdrawsSum = v
	default:
	}

	return withdrawsSum, nil
}

func prepareAddWithdrawalsQuery(order string, user string, sum float64) *query {
	return &query{
		request: addWithdrawals,
		args: []interface{}{
			order,
			user,
			sum,
			time.Now().Format(time.RFC3339),
		},
	}
}

func prepareWithdrawalsSumQuery(user string) *query {
	return &query{
		request: getUserWithdrawalsTotalSum,
		args: []interface{}{
			user,
		},
	}
}

func prepareGetAllUserWithdrawals(user string) *query {
	return &query{
		request: getUserWithdrawals,
		args: []interface{}{
			user,
		},
	}
}
