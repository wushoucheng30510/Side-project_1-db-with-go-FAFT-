package main

import (
	"backend/database/dbInfo"
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUserName = "kevinwu"
	dbPassword = "xxxxxxxx"
	dbIp       = "xxxxxxxx"
	dbName     = "FAFT_TEST"
	tableName  = "Result"
)

func main() {
	dbInfo.LogPrint("Parse json file")
	jsonFile := "../../tmp/result.json"
	content, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	defer func() {
		if err := os.Remove(jsonFile); err != nil {
			log.Fatalf("Failed to remove %s: %v", jsonFile, err)
		}
	}()

	var resultData dbInfo.Data
	err = json.Unmarshal(content, &resultData)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUserName, dbPassword, dbIp, dbName))
	if err != nil {
		log.Fatal("Failed to open db: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	dbInfo.LogPrint("Database is connected")

	if err := dbInfo.CreateTable(db, tableName); err != nil {
		log.Fatal("Failed to create table: ", err)
	}

	dbInfo.LogPrint("Preprocess data")
	failedReason, errLogPreProcess := logPreprocess(resultData.Result)
	if errLogPreProcess != nil {
		log.Fatal("Failed to preprocess data: ", err)
	}

	dbInfo.LogPrint("Insert data")
	if err := dbInfo.InsertData(db, tableName, resultData.Time, resultData.Tester, resultData.Name, resultData.Board, resultData.Model, resultData.Version, resultData.LogPath, resultData.Result, failedReason); err != nil {
		log.Fatal("Failed to insert data to db: ", err)
	}
}

func logPreprocess(result string) (string, error) {
	failedReason := ""
	finalReason := ""
	if result != "Pass" {
		f, err := os.Open("failReason.txt")
		if err != nil {
			return "", err
		}
		defer func() {
			f.Close()
			os.Remove("failReason.txt")
		}()

		stat, _ := f.Stat()
		reader := bufio.NewReader(f)
		buf := make([]byte, stat.Size())

		for {
			_, err := reader.Read(buf)
			if err != nil {
				if err != io.EOF {
					dbInfo.LogPrint("Failed to read data from log.txt")
				}
				break
			}

		}
		failedReason = string(buf)
		var revisedFailedReason []string
		revisedFailedReason = strings.Split(failedReason, "[0m")
		finalReason = revisedFailedReason[len(revisedFailedReason)-1]
	}
	return finalReason, nil
}
