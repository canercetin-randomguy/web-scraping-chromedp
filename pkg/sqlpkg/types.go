package sqlpkg

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

// Dnnt forget to call SqlConn.CloseConn when you are done with the connection.
type SqlConn struct {
	DB *sql.DB
}
type ClientFileInfo struct {
	Username      string `json:"username"`
	FileExtension string `json:"file_extension"`
	Filepath      string `json:"filepath"`
	CreatedAt     string `json:"created_at"`
	MainLink      string `json:"mainlink"`
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

// USED FOR RANDOM STRING GENERATION.
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())
