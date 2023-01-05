package sqlpkg

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

func (conn *SqlConn) CreateClientSchema() error {
	_, err := conn.DB.Query("CREATE SCHEMA IF NOT EXISTS 'clients'")
	if err != nil {
		return err
	}
	// see if we created it
	res, err := conn.DB.Exec("SHOW SCHEMAS LIKE 'clients'")
	if err != nil {
		return err
	}
	// if there is 0 rows, then we didn't create it
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("clients schema not created")
	}
	// if there is more than 0 rows, but no error, then we created it.
	if err != nil {
		return err
	}
	return nil
}

func (conn *SqlConn) CreateTableSchema() error {
	query, err := conn.DB.Prepare(`
	create table clients.client_info
(
    username   VARCHAR(25) not null,
    password   BINARY(60)  not null,
    email      VARCHAR(60) not null,
    created_at TIMESTAMP   null
);

`)
	if err != nil {
		return err
	}
	res, err := query.Exec()
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	// if there is 0 rows, then we didn't create it
	if rows == 0 {
		return errors.New("clients info table under the clients schema not created")
	}
	return nil
}
