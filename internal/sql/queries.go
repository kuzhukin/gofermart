package sql

const (
	createUsersTableQuery = `CREATE TABLE IF NOT EXISTS users (
		login 		text		NOT NULL,
		token 		text		NOT NULL,
		balance 	bigint		DEFAULT 0,
		PRIMARY KEY ( login )
	);`

	createUserQuery = `INSERT INTO users (login, token) VALUES ($1, $2);`

	// const getUserBalanceQuery = `SELECT balance FROM users WHERE login = $1;`
	getUser        = `SELECT * FROM users WHERE login = $1;`
	getUserByToken = `SELECT * FROM users WHERE token = $1;`

	increaseUserBalanceQuery = `UPDATE users SET balance = balance + $1 WHERE login = $2;`
	decreaseUserBalanceQuery = `UPDATE users SET balance = balance - $1 WHERE login = $2;`

	createOrdersTableQuery = `CREATE TABLE IF NOT EXISTS orders (
		"id" 			text					NOT NULL,
		"status"		text					NOT NULL,
		"accrual"		double precision		NOT NULL,
		"user"			text					NOT NULL,
		"upload_time"	timestamp 				NOT NULL,
		PRIMARY KEY ("id"),
		CHECK ( "status" IN ( 'NEW', 'PROCESSING', 'INVALID', 'PROCESSED') ),
		FOREIGN KEY ( "user" ) REFERENCES users ( "login" ) ON DELETE CASCADE
	);`

	createOrderQuery = `INSERT INTO orders ("id", "user", "status", "accrual", "upload_time") VALUES ($1, $2, 'NEW', 0, current_timestamp);`

	updateOrderAccrualQuery = `UPDATE orders SET status = $1, accrual = $2 WHERE id = $3;`

	getOrderQuery = `SELECT * FROM orders WHERE "id" = $1;`

	getAllOrdersQuery        = `SELECT * FROM orders WHERE "user" = $1 ORDER BY upload_time;`
	getUnexecutedOrdersQuery = `SELECT * FROM orders WHERE "status" IN ('NEW', 'PROCESSING');`

	createWithdrawalsTableQuery = `CREATE TABLE IF NOT EXISTS withdrawals (
		"order"		text				NOT NULL,
		"user"		text				NOT NULL,
		"sum"		double precision	NOT NULL,
		PRIMARY KEY ( "order" )
	);`
)
