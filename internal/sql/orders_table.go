package sql

import (
	"database/sql"
	"time"
)

const (

	// ORDERS TABLE
	createOrdersTableQuery = `CREATE TABLE IF NOT EXISTS orders (
		"id" 			text					NOT NULL,
		"status"		text					NOT NULL,
		"accrual"		double precision		NOT NULL,
		"user"			text					NOT NULL,
		"upload_time"	text 					NOT NULL,
		PRIMARY KEY ("id"),
		CHECK ( "status" IN ( 'NEW', 'PROCESSING', 'INVALID', 'PROCESSED') ),
		FOREIGN KEY ( "user" ) REFERENCES users ( "login" ) ON DELETE CASCADE
	);`

	createOrderQuery = `INSERT INTO orders ("id", "user", "status", "accrual", "upload_time") VALUES ($1, $2, 'NEW', 0, $3);`

	updateOrderAccrualQuery = `UPDATE orders SET status = $1, accrual = $2 WHERE id = $3;`

	getOrderQuery = `SELECT * FROM orders WHERE "id" = $1;`

	getAllOrdersQuery        = `SELECT * FROM orders WHERE "user" = $1 ORDER BY upload_time;`
	getUnexecutedOrdersQuery = `SELECT * FROM orders WHERE "status" IN ('NEW', 'PROCESSING');`
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID           string      `json:"id"`
	User         string      `json:"user"`
	Status       OrderStatus `json:"status"`
	Accrual      float32     `json:"accrual,omitempty"`
	UpdaloadTime string      `json:"uploaded_at"`
}

func (o *Order) scan(rows *sql.Rows) error {
	return rows.Scan(&o.ID, &o.Status, &o.Accrual, &o.User, &o.UpdaloadTime)
}

func prepareCreateOrderQuery(orderID string, user string) *query {
	return &query{
		request: createOrderQuery,
		args: []interface{}{
			orderID,
			user,
			time.Now().Format(time.RFC3339),
		},
	}
}

func prepareGetOrderQuery(orderID string) *query {
	return &query{
		request: getOrderQuery,
		args: []interface{}{
			orderID,
		},
	}
}

func prepareGetAllOrdersQuery(user string) *query {
	return &query{
		request: getAllOrdersQuery,
		args: []interface{}{
			user,
		},
	}
}

func prepareGetUnexecutedOrdersQuery() *query {
	return &query{
		request: getUnexecutedOrdersQuery,
		args:    []interface{}{},
	}
}
