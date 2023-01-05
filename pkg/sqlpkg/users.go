package sqlpkg

import (
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

func (conn *SqlConn) GetExistingUser(database string) error {
	// results, err := conn.DB.Query("SELECT * FROM clients.client_info")
	return nil
}
