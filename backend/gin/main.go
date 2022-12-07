package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var dataResult []result

const (
	dbUserName = "test"
	dbPassword = "Test123456"
	dbIp       = "10.240.102.12"
	dbName     = "FAFT_test"
	tableName  = "Result"
)

type result struct {
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
	Log        string `json:"log"`
}

func main() {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUserName, dbPassword, dbIp, dbName))
	if err != nil {
		log.Fatal("Failed to open db: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	log.Print("Database is connected")

	err2 := SearhData(db)
	if err2 != nil {
		log.Fatal("Failed to search data in table `Result`: ", err)
	}

	router.GET("/test", func(ctx *gin.Context) {
		err2 := SearhData(db)
		if err2 != nil {
			log.Fatal("Failed to search data in table `Result`: ", err)
		}
		ctx.IndentedJSON(http.StatusOK, dataResult)
		dataResult = nil
	})

	router.GET("/test2", func(ctx *gin.Context) {
		ctx.FileAttachment("./log/test.txt", "try.txt")
	})

	router.Run(":8082")
}

func SearhData(db *sql.DB) error {
	query := fmt.Sprint("select id,time,tester,name,board,model,version,logPath,result,reason,log from Result order by id desc")
	rows, err := db.Query(query)
	if err != nil {
		return err
	}

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
			log        string
		)
		rows.Scan(&id, &time, &tester, &name, &board, &model, &version, &logPath, &passOrFail, &reason, &log)

		data := result{
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
			Log:        log,
		}
		dataResult = append(dataResult, data)
	}
	return nil
}
