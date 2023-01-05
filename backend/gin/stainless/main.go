package main

import (
	"backend/database/dbInfo"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUserName = "cienet"
	dbPassword = "CienetFAFT"
	dbIp       = "localhost"
	dbName     = "FAFT_TEST"
	tableName  = "Stainless_Result"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUserName, dbPassword, dbIp, dbName))
	if err != nil {
		log.Fatal("Failed to open db: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	log.Print("Database is connected")

	// Create Stainless_Result if it's not existed.
	if err := dbInfo.CreateStainlessTable(db, "Stainless_Result"); err != nil {
		log.Fatal("Failed to create stainless table in DB: ", err)
	}

	// Read csv files and insert data to DB
	// Right now the data should have the following information
	// time, duration, suite, board, model, build version, host, test name, status, reason, firmware RO version, firmware RW version
	dbInfo.LogPrint("Inserting data to DB")
	rows := readCSV(os.Args[1])
	for i, data := range rows {
		if i == 0 {
			continue
		}
		errInsert := dbInfo.InsertStainlessData(db, tableName, data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7], data[8], data[9], data[10], data[11])
		if errInsert != nil {
			dbInfo.LogPrint(fmt.Sprintf("Failed to insert data into DB: %v", errInsert))
		}
	}
	dbInfo.LogPrint("Insert Successfully")
}

// readCSV would read a csv file and put it in a slice
func readCSV(csvFile string) [][]string {
	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open %q: %v", csvFile, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ','

	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal("Cannot read CSV data:", err)
	}
	return rows
}
