package db

import "database/sql"

func Query() (q *Queries, err error) {
	var conn *sql.DB
	conn, err = Connect()
	if err != nil {
		return
	}

	q = New(conn)

	return
}

func QueryWithTx() (tx *sql.Tx, q *Queries, err error) {
	var conn *sql.DB
	conn, err = Connect()
	if err != nil {
		return
	}

	tx, err = conn.Begin()
	if err != nil {
		return
	}

	q = &Queries{db: tx}

	return
}
