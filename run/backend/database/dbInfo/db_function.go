package dbInfo

import (
	"database/sql"
	"fmt"
)

type Result struct {
	Id         int    `json:"id"`
	Time       string `json:"time"`
	Tester     string `json:"tester"`
	Name       string `json:"name"`
	Board      string `json:"board"`
	Model      string `json:"model"`
	Version    string `json:"version"`
	LogPath    string `json:"logPath"`
	PassOrFail string `json:"passOrFail"`
	Reason     string `json:"reason"`
}

func CreateTable(DB *sql.DB, tableName string) error {
	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
            id int AUTO_INCREMENT PRIMARY KEY,
            time datetime,
            tester varchar(15),
            name varchar(30),
            board varchar(20),
            model varchar(20),
            version varchar(20),
            logPath varchar(40),
            result varchar(5),
            reason text
        ); `, tableName)

	if _, err := DB.Exec(sql); err != nil {
		return err
	}
	LogPrint("Table is created")
	return nil
}

func InsertData(DB *sql.DB, tableName, time, tester, name, board, model, version, logPath, result, reason string) error {
	sqlStmt, err := DB.Prepare(fmt.Sprintf("INSERT %s SET time=?,tester=?,name=?,board=?,model=?,version=?,logPath=?,result=?,reason=?", tableName))
	if err != nil {
		return err
	}

	if _, err := sqlStmt.Exec(time, tester, name, board, model, version, logPath, result, reason); err != nil {
		return err
	}
	LogPrint("Insert successfully")
	return nil
}
