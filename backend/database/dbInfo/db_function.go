package dbInfo

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
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

func CreateStainlessTable(DB *sql.DB, tableName string) error {
	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
			id int AUTO_INCREMENT PRIMARY KEY,
			time TIMESTAMP,
			duration varchar(10),
			suite varchar(30),
			board varchar(20),
			model varchar(20), 
			buildVersion varchar(20),
			host varchar(40),
			testName varchar(60),
			status varchar(10),
			reason blob,
			firmwareROVersion varchar(50),
			firmwareRWVersion varchar(50)
        ); `, tableName)

	if _, err := DB.Exec(sql); err != nil {
		return err
	}
	LogPrint("Table is created")
	return nil
}

func InsertStainlessData(DB *sql.DB, tableName, time, duration, suite, board, model, buildVersion, host, testName, status, reason, firmwareROVersion, firmwareRWVersion string) error {
	sqlStmt, err := DB.Prepare(fmt.Sprintf("INSERT %s SET time=?,duration=?,suite=?,board=?,model=?,buildVersion=?,host=?,testName=?,status=?,reason=?,firmwareROVersion=?,firmwareRWVersion=?", tableName))
	if err != nil {
		return err
	}

	if _, err := sqlStmt.Exec(time, duration, suite, board, model, buildVersion, host, testName, status, reason, firmwareROVersion, firmwareRWVersion); err != nil {
		return err
	}
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

func ValidCsv(inputCsv string) error {
	csvD1, csvD2, err := checkCsvInputFormat(inputCsv)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while checking file: %v", err))
	}

	cmdString := "ls /home/ubuntu/backend/gin/stainless | grep .csv"
	cmd := exec.Command("bash", "-c", cmdString)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	LogPrint("Searching the logs in Server...")
	if err := cmd.Start(); err != nil {
		return errors.New(fmt.Sprintf("while running ls command: %v", err))
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return errors.New(fmt.Sprintf("while reading the command output: %v", err))
	}

	if err := cmd.Wait(); err != nil {
		// "exit status 1" means "not found" here.
		if !strings.Contains(err.Error(), "exit status 1") {
			log.Fatal("Unexpected error happened while running ls command: ", err)
		}
	} else {
		if err := verifyOverlap(bytes, csvD1, csvD2); err != nil {
			return errors.New(fmt.Sprintf("while verifying duplicated log in DB: %v", err))
		}
	}
	LogPrint(fmt.Sprintf("Date: %v", string(bytes)))
	LogPrint("Valid log name")

	return nil
}

// checkCsvInputFormat verifies the input file name should be `\d+-\d+.csv`
// Rules:
// 		1. 8 digits date [$year$month$day]
//      2. Year: after 2020
//      3. Month: [1-12]
//      4. Day: [1-31]
func checkCsvInputFormat(inputCsv string) (int, int, error) {
	// Confirm the format is correct
	nameReg, _ := regexp.Compile(`\d+-\d+.csv`)
	matches := nameReg.FindStringSubmatch(inputCsv)
	if len(matches) < 1 {
		return 0, 0, errors.New(fmt.Sprint("csv naming format error: it should be [nums]-[nums].csv"))
	}

	csvName := strings.Split(inputCsv, ".csv")
	csVDate := strings.Split(csvName[0], "-")

	// Confirm it is 8 digits
	if len(csVDate[0]) != 8 || len(csVDate[1]) != 8 {
		return 0, 0, errors.New(fmt.Sprint("csv naming format error: the date should be 8 digits"))
	}

	// Confirm the name before - can coverted to be a integer
	csvDate1, err := strconv.Atoi(csVDate[0])
	if err != nil {
		return 0, 0, errors.New(fmt.Sprintf("csv naming format error: %s is not an integer", csvName[0]))
	}

	// Confirm the name before - can coverted to be a integer
	csvDate2, err := strconv.Atoi(csVDate[1])
	if err != nil {
		return 0, 0, errors.New(fmt.Sprintf("csv naming format error: %s is not an integer", csvName[1]))
	}

	//Confirm year format
	yearMin := 2020
	if (csvDate1/10000) < yearMin || (csvDate2/10000) < yearMin {
		return 0, 0, errors.New(fmt.Sprintf("csv naming format error: year should be greater than %d", yearMin))
	}

	if (csvDate1%10000)/100 > 13 || (csvDate1%10000)/100 < 1 || (csvDate2%10000)/100 > 13 || (csvDate2%10000)/100 < 1 {
		return 0, 0, errors.New(fmt.Sprint("csv naming format error: month should be [1,12]"))
	}
	// Confirm month format
	if (csvDate1%100) > 31 || (csvDate1%100) < 1 || (csvDate2%100) > 31 || (csvDate2%100) < 1 {
		return 0, 0, errors.New(fmt.Sprint("csv naming format error: date should be [1,31]"))
	}
	//Confirm date format

	if csvDate1 >= csvDate2 {
		return 0, 0, errors.New(fmt.Sprintf("csv naming format error: %d should less than %d", csvDate1, csvDate2))
	}
	return csvDate1, csvDate2, nil
}

// verifyOverlap would verify the same period of data would not be duplicated.
// It would be verified by the name of csv file.
func verifyOverlap(bytes []byte, csvD1, csvD2 int) error {
	dateIntSlice := []int{}
	stainlessFolder := strings.Split(string(bytes), "\n")

	for index, name := range stainlessFolder {
		if index == len(stainlessFolder)-1 {
			break
		}

		csvNameSlice := strings.Split(name, "-")

		date1, err := strconv.Atoi(csvNameSlice[0])
		if err != nil {
			log.Println("Noise data is imported. Check the stainless folder manually")
			continue
		}

		date2, err2 := strconv.Atoi(csvNameSlice[1][0 : len(csvNameSlice[1])-4])
		if err2 != nil {
			log.Println("Noise data is imported. Check the stainless folder manually")
			continue
		}
		dateIntSlice = append(dateIntSlice, date1)
		dateIntSlice = append(dateIntSlice, date2)
	}

	sortDateSlice := sort.IntSlice(dateIntSlice)
	sort.Sort(sortDateSlice)

	log.Println("Current log: ", sortDateSlice)
	csvD1Index := sort.Search(len(sortDateSlice), func(i int) bool {
		return sortDateSlice[i] > csvD1
	})

	csvD2Index := sort.Search(len(sortDateSlice), func(i int) bool {
		return sortDateSlice[i] > csvD2
	})

	if sortDateSlice[0] == csvD1 {
		return errors.New(fmt.Sprintf("[%d] has appeared in stainless folder", csvD1))
	} else if csvD1Index-1 >= 0 && sortDateSlice[csvD1Index-1] == csvD1 {
		return errors.New(fmt.Sprintf("[%d] has appeared in stainless folder", csvD1))
	}
	if sortDateSlice[0] == csvD2 {
		return errors.New(fmt.Sprintf("[%d] has appeared in stainless folder", csvD2))
	} else if csvD2Index-1 >= 0 && sortDateSlice[csvD2Index-1] == csvD2 {
		return errors.New(fmt.Sprintf("[%d] has appeared in stainless folder", csvD2))
	}

	if csvD1Index != csvD2Index {
		if csvD1Index%2 == 0 {
			return errors.New(fmt.Sprintf("[ %d-%d ] has overlapped the data", sortDateSlice[csvD1Index], sortDateSlice[csvD1Index+1]))
		} else {
			return errors.New(fmt.Sprintf("[ %d-%d ] has overlapped the data", sortDateSlice[csvD1Index-1], sortDateSlice[csvD1Index]))
		}

	} else {
		if csvD1Index%2 == 1 {
			return errors.New(fmt.Sprintf("[ %d-%d ] has overlapped the data", sortDateSlice[csvD1Index-1], sortDateSlice[csvD1Index]))
		}
	}
	return nil
}
