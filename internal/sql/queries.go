package sql

const createUsersTableQuery = `CREATE TABLE IF NOT EXIST users (
	login 		text		NOT NULL,
	token 		text		NOT NULL,
	balance 	bigint		DEFAULT 0,
	PRIMARY KEY ( login )
);`

const addNewUserQuery = `INSERT INTO users (login, token) VALUES ($1, $2);`
const getUserBalanceQuery = `SELECT balance FROM users WHERE login = $1;`
const getUserToken = `SELECT token FROM users WHERE login = $1;`
const increaseUserBalanceQuery = `UPDATE users SET balance = balance + $1 WHERE login = $2;`
const decreaseUserBalanceQuery = `UPDATE users SET balance = balance - $1 WHERE login = $2;`

const createOrdersTableQuery = `CREATE TABLE IF NOT EXIST orders (
	id 			text	NOT NULL,
	status		text	NOT NULL,
	user		text	DEFAULT 'NEW',
	PRIMARY KEY ( id ),
	CHECK ( 
		status IN ( 'NEW', 'PROCESSING', 'INVALID', 'PROCESSED')
	),
	FOREIGN KEY ( user ) REFERENCES users ( login ) ON DELETE CASCADE,
);`

const addNewOrderQuery = `INSERT INTO orders (id, user) VALUES ($1, $2);`
const updateOrderStatusQuery = `UPDATE orders SET status = $1 WHERE id = $2;`
const getUserOrders = `SELECT id, status FROM orders WHERE user = $1`
