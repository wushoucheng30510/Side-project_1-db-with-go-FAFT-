package main

import (
	"backend/database/dbInfo"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// var dataResult []result

const (
	dbUserName = "kevinwu"
	dbPassword = "CienetFAFT"
	dbIp       = "10.240.102.16"
	dbName     = "FAFT_TEST"
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
}

var dataResult []dbInfo.Result

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

	// Search data for testing results
	var err2 error
	router.GET("/test", func(ctx *gin.Context) {
		dataResult, err2 = dbInfo.SearhData(db)
		if err2 != nil {
			log.Fatal("Failed to search data in table `Result`: ", err)
		}
		ctx.IndentedJSON(http.StatusOK, dataResult)
		dataResult = nil
	})

	// Upload csv files
	router.POST("/uploadCSV", func(c *gin.Context) {
		// get file from form input name 'file'
		file, _ := c.FormFile("stainlessData")

		log.Println(file.Filename)
		if err := dbInfo.ValidCsv(file.Filename); err != nil {
			c.String(http.StatusOK, "Error happend: ", err)
		} else {
			c.SaveUploadedFile(file, "stainless/"+file.Filename)
			c.String(http.StatusOK, "file: %s", file.Filename)
		}

	})

	router.Static("/logDB", "/home/ubuntu/backend/logDB")
	router.Run(":8082")
}
