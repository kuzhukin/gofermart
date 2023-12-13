package sql

const createUsersTableQuery = `CREATE TABLE IF NOT EXISTS users (
	login 		text		NOT NULL,
	token 		text		NOT NULL,
	balance 	bigint		DEFAULT 0,
	PRIMARY KEY ( login )
);`

const createUserQuery = `INSERT INTO users (login, token) VALUES ($1, $2);`
const getUserBalanceQuery = `SELECT balance FROM users WHERE login = $1;`
const getUser = `SELECT * FROM users WHERE login = $1;`
const getUserByToken = `SELECT * FROM users WHERE token = $1;`
const increaseUserBalanceQuery = `UPDATE users SET balance = balance + $1 WHERE login = $2;`
const decreaseUserBalanceQuery = `UPDATE users SET balance = balance - $1 WHERE login = $2;`

const createOrdersTableQuery = `CREATE TABLE IF NOT EXISTS orders (
	"id" 			text	NOT NULL,
	"status"		text	NOT NULL DEFAULT 'NEW',
	"user"			text	NOT NULL,
	PRIMARY KEY ("id"),
	CHECK ( "status" IN ( 'NEW', 'PROCESSING', 'INVALID', 'PROCESSED') ),
	FOREIGN KEY ( "user" ) REFERENCES users ( "login" ) ON DELETE CASCADE
);`

const createOrderQuery = `INSERT INTO orders (id, user) VALUES ($1, $2);`
const updateOrderStatusQuery = `UPDATE orders SET status = $1 WHERE id = $2;`
const getUserOrdersQuery = `SELECT id, status FROM orders WHERE user = $1;`
const getOrderQuery = `SELECT * FROM orders WHERE id = $1 AND user = $2;`

const createWithdrawalsTableQuery = `CREATE TABLE IF NOT EXISTS withdrawals (
	"order"		text	NOT NULL,
	"user"		text	NOT NULL,
	"sum"		bigint	NOT NULL,
	PRIMARY KEY ( "order" )
);`
