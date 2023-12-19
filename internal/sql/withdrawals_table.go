package sql

import (
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

func ScanNewWithdrawRecord(scanner objectScanner) (*UserWithdrawRecord, error) {
	u := &UserWithdrawRecord{}
	err := scanner.Scan(&u.OrderID, &u.Accrual, &u.ProcessedAt)
	if err != nil {
		return nil, fmt.Errorf("withdraw record scan err=%w", err)
	}

	return u, nil
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
