package sqlpkg

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func (conn *SqlConn) CreateNewUser(username string, password []byte, email string, created_at string) error {
	_, err := conn.DB.Query(
		"INSERT INTO clients.client_info (username, password, email, created_at) VALUES (?, ?, ?, ?)",
		username, password, email, created_at)
	if err != nil {
		return err
	}
	return nil
}

func (conn *SqlConn) GetExistingUserPassword(username string) (error, []byte) {
	var password []byte
	results, err := conn.DB.Query(`
	SELECT password
	FROM clients.client_info
	WHERE username = ?`, username)
	if err != nil {
		return err, []byte("")
	}
	defer func(results *sql.Rows) {
		err = results.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(results)
	for results.Next() {
		err = results.Scan(&password)
		if err != nil {
			return err, []byte("")
		}
	}
	return nil, password
}
