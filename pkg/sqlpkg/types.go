package sqlpkg

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Dnnt forget to call SqlConn.CloseConn when you are done with the connection.
type SqlConn struct {
	DB *sql.DB
}

var RunningOnDocker = false

// these are self-explanatory, change them to your own database credentials when you are running them locally.
//
// will change it to work with docker-composer environment variables later.
var (
	Username = "cansu"
	Password = "1234"
	Hostname = "(127.0.0.1:3306)"
	Schema   = "clients"
	Table    = "client_info"
)
