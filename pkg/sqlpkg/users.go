package sqlpkg

import (
	_ "github.com/go-sql-driver/mysql"
	"strings"
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
	for results.Next() {
		err = results.Scan(&password)
		if err != nil {
			return err, []byte("")
		}
	}
	err = results.Close()
	if err != nil {
		return err, []byte("")
	}
	return nil, password
}

// Thanks so much https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go !
func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func (conn *SqlConn) InsertAuthenticationToken(username string, auth string) error {
	_, err := conn.DB.Query("UPDATE clients.client_info SET auth_token = ? WHERE username = ?", auth, username)
	if err != nil {
		return err
	}
	return nil
}
func (conn *SqlConn) RetrieveAuthenticationToken(username string) (string, error) {
	var auth string
	results, err := conn.DB.Query("SELECT auth_token FROM clients.client_info WHERE username = ?", username)
	if err != nil {
		return "", err
	}
	for results.Next() {
		err = results.Scan(&auth)
		if err != nil {
			return "", err
		}
	}
	err = results.Close()
	if err != nil {
		return "", err
	}
	return auth, nil
}

func (conn *SqlConn) InsertUserPackageDetails(username string, planName string) error {
	_, err := conn.DB.Query("UPDATE clients.client_info SET package = ? WHERE username = ?", planName, username)
	if err != nil {
		return err
	}
	return nil
}
func (conn *SqlConn) RetrieveUserPackageDetails(username string) (string, error) {
	var pkg string
	results, err := conn.DB.Query("SELECT package FROM clients.client_info WHERE username = ?", username)
	if err != nil {
		return "", err
	}
	for results.Next() {
		err = results.Scan(&pkg)
		if err != nil {
			return "", err
		}
	}
	err = results.Close()
	if err != nil {
		return "", err
	}
	return pkg, nil
}

func (conn *SqlConn) InsertUserLinkLimit(username string, limit int) error {
	_, err := conn.DB.Query("UPDATE clients.client_info SET link_limit = ? WHERE username = ?", limit, username)
	if err != nil {
		return err
	}
	return nil
}

func (conn *SqlConn) RetrieveUserLinkLimit(username string) (int, error) {
	var limit int
	results, err := conn.DB.Query("SELECT link_limit FROM clients.client_info WHERE username = ?", username)
	if err != nil {
		return 0, err
	}
	for results.Next() {
		err = results.Scan(&limit)
		if err != nil {
			return 0, err
		}
	}
	err = results.Close()
	if err != nil {
		return 0, err
	}
	return limit, nil
}

func (conn *SqlConn) InsertFileLink(username string, link string, created_at string, fileExtension string) error {
	_, err := conn.DB.Query("INSERT INTO clients.client_file_info (username, file_extension,filepath,created_at) VALUES (?, ?, ?,?)", username, fileExtension, link, created_at)
	if err != nil {
		return err
	}
	return nil
}
func (conn *SqlConn) RetrieveFileLinks(username string, fileExtension string) ([]ClientFileInfo, error) {
	rows, err := conn.DB.Query("SELECT username, file_extension, filepath, created_at FROM clients.client_file_info WHERE username = ? AND file_extension = ?", username, fileExtension)
	if err != nil {
		return nil, err
	}
	var files []ClientFileInfo
	for rows.Next() {
		var file ClientFileInfo
		err = rows.Scan(&file.Username, &file.FileExtension, &file.Filepath, &file.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}
