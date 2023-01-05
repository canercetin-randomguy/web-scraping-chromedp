package sqlpkg

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// CloseConn closes the sql connection.
func (conn *SqlConn) CloseConn() error {
	err := conn.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

// GetSQLConn attaches a new sql connection to the SqlConn struct.
func (conn *SqlConn) GetSQLConn(dbname string) error {
	db, err := sql.Open("mysql", Username+":"+Password+"@tcp"+Hostname+"/"+dbname+"?parseTime=true")
	if err != nil {
		log.Println("Error", err.Error())
		return err
	}
	conn.DB = db
	return nil
}
