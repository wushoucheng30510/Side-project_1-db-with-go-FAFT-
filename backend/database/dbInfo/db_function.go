package dbInfo

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
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

type existedLog struct {
	LogPath string `json:"logPath"`
}

// var logOutput []logResult

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
	log.Println("Table is created")
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
	log.Println("Insert successfully")
	return nil
}

func SearhData(db *sql.DB) ([]Result, error) {
	query := fmt.Sprint("select id,time,tester,name,board,model,version,logPath,result,reason from Result order by id desc")
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var dataResult []Result
	for rows.Next() {
		var (
			id         int
			time       string
			tester     string
			name       string
			board      string
			model      string
			version    string
			logPath    string
			passOrFail string
			reason     string
		)
		rows.Scan(&id, &time, &tester, &name, &board, &model, &version, &logPath, &passOrFail, &reason)

		data := Result{
			Id:         id,
			Time:       time,
			Tester:     tester,
			Name:       name,
			Board:      board,
			Model:      model,
			Version:    version,
			LogPath:    logPath,
			PassOrFail: passOrFail,
			Reason:     reason,
		}
		dataResult = append(dataResult, data)
	}
	return dataResult, nil
}

func CleanLog(DB *sql.DB) error {
	LogPrint("Start to clean up the log in Server")

	logInLogDBFolder, errLogInServer := searchExistedLogInServer()
	if errLogInServer != nil {
		return errors.New(fmt.Sprint("failed to search the logs in Server: ", errLogInServer))
	}

	logInDB, err := searchExistedLogInDB(DB)
	if err != nil {
		return errors.New(fmt.Sprint("failed to search log in DB: ", err))
	}

	var notNeededLog []string
	var logToRemove string
	notNeeded := 0
	for _, logInFolder := range logInLogDBFolder {
		for index, logDB := range logInDB {
			exitedLog := strings.Contains(logInFolder, logDB.LogPath)
			if exitedLog {
				break
			} else {
				if index == len(logInDB)-1 {
					notNeeded += 1
					notNeededLog = append(notNeededLog, logInFolder)
					if len(logInFolder) != 0 {
						logToRemove = logToRemove + "../database/logDB/" + logInFolder + " "
					}
				} else {
					continue
				}
			}
		}
	}
	log.Print("Not use log: ", notNeeded)
	if err := removeLog(logToRemove); err != nil {
		return errors.New(fmt.Sprint("failed to remove logs: ", err))
	}
	return nil
}

func searchExistedLogInServer() ([]string, error) {
	cmdString := "ls ../database/logDB"
	cmd := exec.Command("bash", "-c", cmdString)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	LogPrint("Searching the logs in Server...")
	if err := cmd.Start(); err != nil {
		return nil, errors.New(fmt.Sprint("failed to run the command: ", err))
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, errors.New(fmt.Sprint("failed to read all output: ", err))
	}

	if err := cmd.Wait(); err != nil {
		return nil, errors.New(fmt.Sprint("failed to wait for output be completed: ", err))
	}

	logInLogDBFolder := strings.Split(string(bytes), "\n")
	return logInLogDBFolder, nil
}

func searchExistedLogInDB(DB *sql.DB) ([]existedLog, error) {
	query := fmt.Sprintf("select logPath from Result order by id desc")
	rows, err := DB.Query(query)

	var logOutput []existedLog
	if err != nil {
		return nil, err
	}

	LogPrint("Searching the logs in DB...")
	for rows.Next() {
		var logPath string
		rows.Scan(&logPath)

		data := existedLog{
			LogPath: logPath,
		}
		logOutput = append(logOutput, data)
	}
	return logOutput, nil
}

func removeLog(logToRemove string) error {
	if len(logToRemove) == 0 {
		return nil
	}
	cmdString := fmt.Sprintf("rm %s", logToRemove)
	fmt.Println(cmdString)
	cmd := exec.Command("bash", "-c", cmdString)

	LogPrint("Removing the no used logs ...")
	if err := cmd.Start(); err != nil {
		return errors.New(fmt.Sprint("failed to run the command: ", err))
	}

	if err := cmd.Wait(); err != nil {
		return errors.New(fmt.Sprint("failed to wait for command to be completed: ", err))
	}
	return nil
}
