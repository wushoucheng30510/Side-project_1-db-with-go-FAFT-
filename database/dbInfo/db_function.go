package dbInfo

import (
	"database/sql"
	"fmt"
)

func CreateTable(DB *sql.DB, tableName string) error {
	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS Result_%s(
            id int AUTO_INCREMENT PRIMARY KEY,
            time TIMESTAMP(6),
            tester varchar(15),
            name varchar(30),
            board varchar(20),
            model varchar(20),
            version varchar(20),
            logPath varchar(40),
            result varchar(5),
            reason text,
            log mediumtext
        ); `, tableName)

	if _, err := DB.Exec(sql); err != nil {
		return err
	}
	LogPrint("Table is created")
	return nil
}

func InsertData(DB *sql.DB, tableName, time, tester, name, board, model, version, logPath, result, reason, testlog string) error {
	sqlStmt, err := DB.Prepare(fmt.Sprintf("INSERT Result_%s SET time=?,tester=?,name=?,board=?,model=?,version=?,logPath=?,result=?,reason=?,log=?", tableName))
	if err != nil {
		return err
	}

	if _, err := sqlStmt.Exec(time, tester, name, board, model, version, logPath, result, reason, testlog); err != nil {
		return err
	}
	LogPrint("Insert successfully")
	return nil
}
