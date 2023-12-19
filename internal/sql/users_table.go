package sql

import (
	"errors"
	"fmt"
)

const (
	createUsersTableQuery = `CREATE TABLE IF NOT EXISTS users (
		login 		text					NOT NULL,
		token 		text					NOT NULL,
		balance 	double precision		DEFAULT 0,
		PRIMARY KEY ( login )
	);`

	createUserQuery = `INSERT INTO users (login, token) VALUES ($1, $2);`

	getUser        = `SELECT * FROM users WHERE login = $1;`
	getUserByToken = `SELECT * FROM users WHERE token = $1;`

	increaseUserBalanceQuery = `UPDATE users SET balance = balance + $1 WHERE login = $2;`
	decreaseUserBalanceQuery = `UPDATE users SET balance = balance - $1 WHERE login = $2;`
)

type User struct {
	Login     string
	AuthToken string
	Balance   float64
}

func (u *User) Scan(scanner objectScanner) error {
	return scanner.Scan(&u.Login, &u.AuthToken, &u.Balance)
}

var ErrEmptyScannerResult = errors.New("sql obj scanner has empty result")

func ScanNewUser(scanner objectScanner) (*User, error) {
	if !scanner.Next() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		return nil, ErrEmptyScannerResult
	}

	u := &User{}
	err := u.Scan(scanner)
	if err != nil {
		return nil, fmt.Errorf("user scan err=%w", err)
	}

	return u, nil
}

func prepareCreateUserQuery(login, token string) *query {
	return &query{
		request: createUserQuery,
		args:    []interface{}{login, token},
	}
}

func prepareGetUserQuery(login string) *query {
	return &query{
		request: getUser,
		args:    []interface{}{login},
	}
}

func prepareGetUserByTokenQuery(token string) *query {
	return &query{
		request: getUserByToken,
		args:    []interface{}{token},
	}
}

// func prepareIncreaseUserBalanceQuery(login string, balance float64) *query {
// 	return &query{
// 		request: increaseUserBalanceQuery,
// 		args:    []interface{}{balance, login},
// 	}
// }

func prepareDecreaseUserBalanceQuery(login string, balance float64) *query {
	return &query{
		request: decreaseUserBalanceQuery,
		args:    []interface{}{balance, login},
	}
}
